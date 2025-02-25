//go:build windows

package wirekcfg

import (
	"errors"
	"fmt"
	"net/netip"
	"wirekcp/wgengine/winnet"

	ole "github.com/go-ole/go-ole"
	"github.com/wirekcp/wireguard-go/tun"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
)

func SetIP(tunnel tun.Device, config *Config) error {
	nativeTunDevice, ok := tunnel.(*tun.NativeTun)
	if ok {
		link := winipcfg.LUID(nativeTunDevice.LUID())
		ip, err := netip.ParsePrefix(config.IPv4CIDR)
		if err != nil {
			return err
		}
		err = link.SetIPAddresses([]netip.Prefix{ip})
		if err != nil {
			return err
		}
		_, err = setPrivateNetwork(link)
		if err != nil {
			return err
		}
	} else {
		return errors.New("tunnel is not a NativeTun")
	}
	return nil
}

func setPrivateNetwork(ifcLUID winipcfg.LUID) (bool, error) {
	// NLM_NETWORK_CATEGORY values.
	const (
		categoryPublic  = 0
		categoryPrivate = 1
		categoryDomain  = 2
	)

	ifcGUID, err := ifcLUID.GUID()
	if err != nil {
		return false, fmt.Errorf("ifcLUID.GUID: %v", err)
	}

	// aaron: DO NOT call Initialize() or Uninitialize() on c!
	// We've already handled that process-wide.
	var c ole.Connection

	m, err := winnet.NewNetworkListManager(&c)
	if err != nil {
		return false, fmt.Errorf("winnet.NewNetworkListManager: %v", err)
	}
	defer m.Release()

	cl, err := m.GetNetworkConnections()
	if err != nil {
		return false, fmt.Errorf("m.GetNetworkConnections: %v", err)
	}
	defer cl.Release()

	for _, nco := range cl {
		aid, err := nco.GetAdapterId()
		if err != nil {
			return false, fmt.Errorf("nco.GetAdapterId: %v", err)
		}
		if aid != ifcGUID.String() {
			continue
		}

		n, err := nco.GetNetwork()
		if err != nil {
			return false, fmt.Errorf("GetNetwork: %v", err)
		}
		defer n.Release()

		cat, err := n.GetCategory()
		if err != nil {
			return false, fmt.Errorf("GetCategory: %v", err)
		}

		if cat != categoryPrivate && cat != categoryDomain {
			if err := n.SetCategory(categoryPrivate); err != nil {
				return false, fmt.Errorf("SetCategory: %v", err)
			}
		}
		return true, nil
	}

	return false, nil
}
