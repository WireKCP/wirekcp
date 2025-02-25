//go:build !windows

package wirektun

import (
	"time"

	"github.com/wirekcp/wireguard-go/device"
	"github.com/wirekcp/wireguard-go/tun"
)

// Dummy implementation that does nothing.
func waitInterfaceUp(_ tun.Device, _ time.Duration, _ *device.Logger) error {
	return nil
}
