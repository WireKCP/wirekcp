//go:build !windows

package wirekcfg

import (
	"net"
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
	return SetIPwithTunName(wirektun.DefaultTunName(), cidr)
}

func SetIPwithTunName(name, cidr string) error {
	iface, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	addr, err := netlink.ParseAddr(cidr)
	if err != nil {
		return err
	}
	oldAddrString := GetIPwithTunName(name)
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

func GetIPwithTunName(name string) string {
	ifc, _ := net.InterfaceByName(name)
	addrs, _ := ifc.Addrs()
	for _, a := range addrs {
		if a.(*net.IPNet).IP.To4() != nil {
			return a.String()
		}
	}
	return ""
}
