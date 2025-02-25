//go:build windows

package wgengine

import (
	"net"

	"github.com/wirekcp/wireguard-go/ipc"
)

func UAPIListen(name string) (net.Listener, error) {
	return ipc.UAPIListen(name)
}
