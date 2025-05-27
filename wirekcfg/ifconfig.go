package wirekcfg

import (
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/wirekcp/wgctrl/wgtypes"
)

type Config struct {
	// IPv4CIDR is the CIDR of the IPv4 address.
	IPv4CIDR   string
	ListenPort int
	PrivateKey string
	Peers      []PeerConfig
	Mode       string // "kcp" or "udp"
}

func ReadFromFile(file string) (*Config, error) {
	var config Config
	f, err := os.Open(file)
	if err != nil {
		return nil, os.ErrNotExist
	}
	defer f.Close()
	err = toml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) WriteToFile(file string) error {
	// Check if the file exists overwrite, if not create it
	var f *os.File
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// Create the file if it does not exist
		f, err = os.Create(file)
		if err != nil {
			return err
		}
		defer f.Close()
	} else {
		f, err = os.OpenFile(file, os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	return toml.NewEncoder(f).Encode(c)
}

func (c *Config) ChangeInterface(file string) error {
	oldconfig, err := ReadFromFile(file)
	if err != nil {
		return err
	}

	oldconfig.IPv4CIDR = c.IPv4CIDR
	oldconfig.ListenPort = c.ListenPort
	oldconfig.PrivateKey = c.PrivateKey

	return oldconfig.WriteToFile(file)
}

func (c *Config) AddPeer(peer PeerConfig) error {
	for _, p := range c.Peers {
		if p.PublicKey == peer.PublicKey {
			return os.ErrExist
		} else if p.Endpoint == peer.Endpoint {
			return os.ErrExist
		}
	}
	c.Peers = append(c.Peers, peer)
	return nil
}

func (c *Config) ChangePeer(peer PeerConfig) error {
	for i, p := range c.Peers {
		if p.PublicKey == peer.PublicKey {
			c.Peers[i] = peer
			return nil
		} else if p.Endpoint == peer.Endpoint {
			c.Peers[i] = peer
			return nil
		}
	}
	return os.ErrNotExist
}

func (c *Config) DeletePeer(peer PeerConfig) error {
	for i, p := range c.Peers {
		if p.PublicKey == peer.PublicKey {
			c.Peers = append(c.Peers[:i], c.Peers[i+1:]...)
			return nil
		} else if p.Endpoint == peer.Endpoint {
			c.Peers = append(c.Peers[:i], c.Peers[i+1:]...)
			return nil
		}
	}
	return os.ErrNotExist
}

func (c *Config) ChangePeers(file string) error {
	oldconfig, err := ReadFromFile(file)
	if err != nil {
		return err
	}
	oldconfig.Peers = c.Peers
	if err := oldconfig.WriteToFile(file); err != nil {
		return err
	}
	return nil
}

func (c Config) ToWgConfig() *wgtypes.Config {
	key, _ := wgtypes.ParseKey(c.PrivateKey)
	return &wgtypes.Config{
		PrivateKey: &key,
		ListenPort: &c.ListenPort,
		Peers:      ToPeersConfig(c.Peers),
	}
}
