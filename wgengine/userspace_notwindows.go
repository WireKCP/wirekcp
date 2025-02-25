//go:build !windows

package wgengine

import (
	"net"

	"github.com/wirekcp/wireguard-go/ipc"
)

func UAPIListen(name string) (net.Listener, error) {
	fd, err := ipc.UAPIOpen(name)
	if err != nil {
		return nil, err
	}
	return ipc.UAPIListen(name, fd)
}
