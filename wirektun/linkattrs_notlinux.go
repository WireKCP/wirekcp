//go:build !linux

package wirektun

import "github.com/wirekcp/wireguard-go/tun"

func setLinkAttrs(_ tun.Device) error {
	return nil
}
