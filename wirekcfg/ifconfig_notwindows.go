//go:build !windows

package wirekcfg

import (
	"wirekcp/wirektun"

	"github.com/vishvananda/netlink"
	"github.com/wirekcp/wireguard-go/tun"
)

func SetIP(tunnel tun.Device, config *Config) error {
	name, err := tunnel.Name()
	iface, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	addr, err := netlink.ParseAddr(config.IPv4CIDR)
	if err != nil {
		return err
	}
	err = netlink.AddrAdd(iface, addr)
	if err != nil {
		return err
	}
	return netlink.LinkSetUp(iface)
}

func SetIPwithoutTun(cidr string) error {
	iface, err := netlink.LinkByName(wirektun.DefaultTunName())
	if err != nil {
		return err
	}
	addr, err := netlink.ParseAddr(cidr)
	if err != nil {
		return err
	}
	err = netlink.AddrAdd(iface, addr)
	if err != nil {
		return err
	}
	return netlink.LinkSetUp(iface)
}
