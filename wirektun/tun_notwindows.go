//go:build !windows

package wirektun

import (
	"github.com/wirekcp/wireguard-go/tun"
)

func interfaceName(dev tun.Device) (string, error) {
	return dev.Name()
}
