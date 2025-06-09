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
	var preshared *wgtypes.Key
	if c.Endpoint != "" {
		endPoint, _ = net.ResolveUDPAddr("udp", c.Endpoint)
	} else {
		endPoint = nil
	}
	key, _ := wirektypes.ParseKey(c.PublicKey)
	if c.PresharedKey != "" {
		key, _ := wirektypes.ParseKey(c.PresharedKey)
		preshared = &key
	} else {
		preshared = nil
	}

	return wgtypes.PeerConfig{
		PublicKey:         key,
		PresharedKey:      preshared,
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

func ToPeersConfig(peers []PeerConfig) []wgtypes.PeerConfig {
	wgPeers := make([]wgtypes.PeerConfig, len(peers))
	for i, peer := range peers {
		wgPeers[i] = peer.ToWgPeerConfig()
	}
	return wgPeers
}

func IPNetToString(ipn net.IPNet) string {
	if ipn.IP.To4() != nil {
		return ipn.String()
	}
	return strings.Join(strings.Split(ipn.String(), "/")[:2], "/")
}

func IPsNetToString(ipn []net.IPNet) string {
	var ips []string
	for _, ip := range ipn {
		ips = append(ips, IPNetToString(ip))
	}
	return strings.Join(ips, ",")
}

func ToWkPeersConfig(peers []wgtypes.PeerConfig) []PeerConfig {
	peerConfig := []PeerConfig{}
	for _, p := range peers {
		if p.Remove {
			println("Skipping peer with Remove flag set to true")
			continue
		}
		if p.PresharedKey != nil && p.Endpoint != nil {
			peerConfig = append(peerConfig, PeerConfig{
				PublicKey:    p.PublicKey.String(),
				PresharedKey: p.PresharedKey.String(),
				Endpoint:     p.Endpoint.String(),
				AllowedIPs:   IPsNetToString(p.AllowedIPs),
			})
		} else if p.PresharedKey == nil && p.Endpoint != nil {
			peerConfig = append(peerConfig, PeerConfig{
				PublicKey:  p.PublicKey.String(),
				Endpoint:   p.Endpoint.String(),
				AllowedIPs: IPsNetToString(p.AllowedIPs),
			})
		} else if p.PresharedKey != nil && p.Endpoint == nil {
			peerConfig = append(peerConfig, PeerConfig{
				PublicKey:    p.PublicKey.String(),
				PresharedKey: p.PresharedKey.String(),
				AllowedIPs:   IPsNetToString(p.AllowedIPs),
			})
		} else {
			peerConfig = append(peerConfig, PeerConfig{
				PublicKey:  p.PublicKey.String(),
				AllowedIPs: IPsNetToString(p.AllowedIPs),
			})
		}
	}
	return peerConfig
}
