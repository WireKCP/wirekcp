package wirekcfg

import (
	"net"
	"strings"
	"wirekcp/wirektypes"

	"github.com/wirekcp/wgctrl/wgtypes"
)

type PeerConfig struct {
	// base64 encoded public key
	PublicKey string

	// base64 encoded preshared key
	PresharedKey string

	// Endpoint specifies the remote endpoint to which a device will connect.
	Endpoint string

	AllowedIPs string
}

func (c PeerConfig) ToWgPeerConfig() wgtypes.PeerConfig {
	stringIPs := strings.Split(c.AllowedIPs, ",")
	allowedIPs := make([]net.IPNet, len(stringIPs))
	for i, ip := range stringIPs {
		ip = strings.TrimSpace(ip)
		_, netIPNet, _ := net.ParseCIDR(ip)
		allowedIPs[i] = *netIPNet
	}
	var endPoint *net.UDPAddr
	if c.Endpoint != "" {
		endPoint, _ = net.ResolveUDPAddr("udp", c.Endpoint)
	} else {
		endPoint = nil
	}
	key, _ := wirektypes.ParseKey(c.PublicKey)
	return wgtypes.PeerConfig{
		PublicKey:         key,
		Endpoint:          endPoint,
		ReplaceAllowedIPs: true,
		AllowedIPs:        allowedIPs,
	}
}

func ToDeletePeerConfig(peerName string) wgtypes.PeerConfig {
	key, _ := wgtypes.ParseKey(peerName)
	return wgtypes.PeerConfig{
		Remove:    true,
		PublicKey: key,
	}
}
