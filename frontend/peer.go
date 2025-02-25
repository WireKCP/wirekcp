package frontend

import (
	"bytes"
	"wirekcp/wirekcfg"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/wirekcp/wgctrl/wgtypes"
)

func PeerMenuForm() Selection {
	var selection Selection
	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Selection]().Title("Peer Menu").Options(
				huh.NewOption("Add Peer", Add),
				huh.NewOption("Edit Peer", Edit),
				huh.NewOption("Delete Peer", Delete),
				huh.NewOption("Back", Quit),
			).Value(&selection),
		),
	).WithProgramOptions(tea.WithAltScreen()).Run()
	return selection
}

func PeerForm(wgPeer ...wgtypes.Peer) wgtypes.Config {
	var peer wirekcfg.PeerConfig
	var peers []wgtypes.PeerConfig = make([]wgtypes.PeerConfig, 0)
	if len(wgPeer) > 0 {
		peer.AllowedIPs = func() string {
			var ips string
			for i, ip := range wgPeer[0].AllowedIPs {
				if i > 0 {
					ips += ", "
				}
				ips += ip.String()
			}
			return ips
		}()
		if wgPeer[0].Endpoint != nil {
			peer.Endpoint = wgPeer[0].Endpoint.String()
		}
		emptyKey := wgtypes.Key{}
		if !bytes.Equal(wgPeer[0].PresharedKey[:], emptyKey[:]) {
			peer.PresharedKey = wgPeer[0].PresharedKey.String()
		}
		peer.PublicKey = wgPeer[0].PublicKey.String()
		peers = append(peers, wgtypes.PeerConfig{
			PublicKey: wgPeer[0].PublicKey,
			Remove:    true,
		})
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Endpoint").Placeholder("<IP:PORT>").Value(&peer.Endpoint).Validate(ValidateOptionalUDPAddr),
			huh.NewInput().Title("Remote's Public Key").Placeholder("base64 format").Value(&peer.PublicKey).Validate(ValidateRequiredKey),
			huh.NewInput().Title("Preshared Key [optional]").Placeholder("base64 format").Value(&peer.PresharedKey).Validate(ValidateOptionalKey),
			huh.NewInput().Title("Allowed IPs").Placeholder("CIDR format").Value(&peer.AllowedIPs).Validate(ValidateRequiredCIDRs),
		),
	).WithProgramOptions(tea.WithAltScreen())
	if err := form.Run(); err != nil {
		return wgtypes.Config{}
	}

	peers = append(peers, peer.ToWgPeerConfig())

	return wgtypes.Config{
		Peers: peers,
	}
}

func PeerEditForm(device *wgtypes.Device) wgtypes.Config {
	peer := PeerDeleteForm(device)
	if peer == "" {
		return wgtypes.Config{}
	}
	for _, p := range device.Peers {
		if p.PublicKey.String() == peer {
			return PeerForm(p)
		}
	}
	return wgtypes.Config{}
}

func PeerDeleteForm(device *wgtypes.Device) string {
	var peer string
	options := make([]huh.Option[string], len(device.Peers))
	for i, p := range device.Peers {
		if p.Endpoint == nil {
			label := "Peer: " + "Server for " + p.PublicKey.String()
			options[i] = huh.NewOption(label, p.PublicKey.String())
		} else {
			label := "Peer: " + p.Endpoint.String()
			options[i] = huh.NewOption(label, p.PublicKey.String())
		}
	}
	options = append(options, huh.NewOption("Back", ""))
	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title("Select Peer").Options(
				options...,
			).Value(&peer),
		),
	).WithProgramOptions(tea.WithAltScreen()).Run()
	return peer
}
