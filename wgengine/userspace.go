package wgengine

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"

	// "wirekcp/ipc/namedpipe"
	"wirekcp/wirekcfg"
	"wirekcp/wirektun"
	"wirekcp/wirektypes"

	"github.com/wirekcp/wireguard-go/conn"
	"github.com/wirekcp/wireguard-go/device"
	"github.com/wirekcp/wireguard-go/tun"
)

type EngineConfig struct {
	// Logf is the logging function used by the engine.
	Logger *device.Logger
	// TUN is the tun device used by the engine.
	TUN tun.Device
	// ListenPort is the port on which the engine will listen.
	ListenPort uint16
	// ConfigPath is the path to the WireKCP configuration file.
	configPath string
	// LogPath is the path to the WireKCP log file.
	logPath string
	// EchoRespondToAll determines whether ICMP Echo requests incoming from WireKCP peers
	// will be intercepted and responded to, regardless of the source host.
	EchoRespondToAll bool
}

type userspaceEngine struct {
	logger *device.Logger
	reqCh  chan struct{}
	waitCh chan struct{} // chan is closed when first Close call completes; contrast with closing bool
	tundev tun.Device
	wgdev  *device.Device
	uapi   net.Listener

	// localAddrs is the set of IP addresses assigned to the local
	// tunnel interface. It's used to reflect local packets
	// incorrectly sent to us.
	// localAddrs atomic.Value // of map[packet.IP]bool

	// wgLock sync.Mutex // serializes all wgdev operations; see lock order comment below
	//lastCfgFull         wgcfg.Config
	//lastEngineSigFull   string // of full wireguard config
	//lastEngineSigTrim   string // of trimmed wireguard config

	mu      sync.Mutex // guards following; see lock order comment below
	closing bool       // Close was called (even if we're still closing)
	//statusCallback StatusCallback
	//peerSequence   []wgcfg.Key
	//endpoints      []string
	//pingers        map[wgcfg.Key]*pinger // legacy pingers for pre-discovery peers

	// Lock ordering: wgLock, then mu.
}

func NewUserspaceEngine(logger *device.Logger, tunname string, listenPort uint16, config, logPath string) (Engine, error) {
	if tunname == "" {
		return nil, fmt.Errorf("--tun name must not be blank")
	}

	logger.Verbosef("Starting userspace wireguard engine with tun device %q", tunname)

	tun, name, err := wirektun.New(logger, tunname)
	if err != nil {
		wirektun.Diagnose(logger, tunname, err)
		logger.Verbosef("CreateTUN: %v", err)
		return nil, err
	}
	logger.Verbosef("CreateTUN %q: ok", name)

	conf := EngineConfig{
		Logger:     logger,
		TUN:        tun,
		ListenPort: listenPort,
		configPath: config,
		logPath:    logPath,
	}

	e, err := newUserspaceEngineAdvanced(conf)
	if err != nil {
		return nil, err
	}
	return e, err
}

func newUserspaceEngineAdvanced(conf EngineConfig) (Engine, error) {
	e := &userspaceEngine{
		logger: conf.Logger,
		reqCh:  make(chan struct{}, 1),
		waitCh: make(chan struct{}),
		tundev: conf.TUN,
	}

	// wgdev takes ownership of tundev, will close it when closed.
	e.wgdev = wirekcfg.NewDevice(e.tundev, conn.NewStdNetBind(), e.logger)
	e.wgdev.Up()

	tunname, err := e.tundev.Name()
	if err != nil {
		e.logger.Errorf("Failed to get tun device name: %v", err)
		e.Close()
		return nil, err
	}

	e.logger.Verbosef("Creating WireKCP device %q", tunname)

	uapi, err := UAPIListen(tunname)
	if err != nil {
		e.logger.Errorf("Failed to listen on uapi socket: %v", err)
		e.Close()
		return nil, err
	}
	e.uapi = uapi
	errs := make(chan error)
	go func() {
		for {
			conn, err := uapi.Accept()
			if err != nil {
				errs <- err
				return
			}
			go e.wgdev.IpcHandle(conn)
		}
	}()

	e.logger.Verbosef("UAPI listener started")

	var ifconfig *wirekcfg.Config

	if _, err = os.Stat(conf.configPath); os.IsNotExist(err) {
		// If the Directory does not exist, create it
		configFolderPath := filepath.Dir(conf.configPath)
		if err := os.MkdirAll(configFolderPath, 0755); err != nil {
			e.logger.Errorf("Failed to create config directory: %v", err)
			e.Close()
			return nil, err
		}
		ifconfig = &wirekcfg.Config{
			IPv4CIDR:   "192.168.200.1/24",
			ListenPort: int(conf.ListenPort),
			PrivateKey: wirektypes.GeneratePrivateKey(),
		}
		err = ifconfig.WriteToFile(conf.configPath)
		if err != nil {
			e.logger.Errorf("Failed to write config file: %v", err)
			e.Close()
			return nil, err
		}
	} else {
		ifconfig, err = wirekcfg.ReadFromFile(conf.configPath)
		if err != nil {
			e.logger.Errorf("Failed to read config file: %v", err)
			e.Close()
			return nil, err
		}
	}

	wirekcfg.SetIP(e.tundev, ifconfig)

	wirekcfg.ConfigureDevice(tunname, *ifconfig)

	select {
	case <-errs:
		e.Close()
		return nil, err
	default:
	}

	return e, nil
}

func (e *userspaceEngine) Close() {
	e.mu.Lock()
	if e.closing {
		e.mu.Unlock()
		return
	}
	e.closing = true
	e.mu.Unlock()

	if e.uapi != nil {
		e.uapi.Close()
	}
	e.wgdev.Close()
	e.logger.Verbosef("Shutting down userspace engine")

	// Shut down pingers after tundev is closed (by e.wgdev.Close) so the
	// synchronous close does not get stuck on InjectOutbound.
	close(e.waitCh)
}

func (e *userspaceEngine) Wait() chan struct{} {
	return e.waitCh
}
