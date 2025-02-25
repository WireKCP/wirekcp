//go:build !windows

package wirektun

import (
	"github.com/vishvananda/netlink"
	"github.com/wirekcp/wireguard-go/tun"
)

func interfaceName(dev tun.Device) (string, error) {
	return dev.Name()
}

func SetIPwithoutTun(cidr string) error {
	iface, err := netlink.LinkByName(DefaultTunName())
	if err != nil {
		return err
	}
	addr, err := netlink.ParseAddr(cidr)
	if err != nil {
		return err
	}
	oldAddrString := GetIP()
	if oldAddrString != "" {
		oldAddr, err := netlink.ParseAddr(oldAddrString)
		if err != nil {
			return err
		}
		err = netlink.AddrDel(iface, oldAddr)
		if err != nil {
			return err
		}
	}
	err = netlink.AddrAdd(iface, addr)
	if err != nil {
		return err
	}
	return netlink.LinkSetUp(iface)
}
