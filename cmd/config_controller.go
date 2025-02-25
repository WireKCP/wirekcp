package cmd

import (
	"wirekcp/frontend"
	"wirekcp/wirekcfg"
	"wirekcp/wirektun"

	"github.com/spf13/cobra"
	"github.com/wirekcp/wgctrl"
	"github.com/wirekcp/wgctrl/wgtypes"
)

var (
	interfaceCmd = &cobra.Command{
		Use:   "interface",
		Short: "Interface Configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _ = wgctrl.New()
			clientDevice, _ = client.Device(wirektun.DefaultTunName())
			if err := frontend.ValidateDeviceExists(client, wirektun.DefaultTunName()); err != nil {
				return err
			}
			config := frontend.InterfaceForm(*clientDevice)
			if config == nil {
				return nil
			}
			wirektun.SetIPwithoutTun(config.IPv4CIDR)
			return client.ConfigureDevice(clientDevice.Name, *config.ToWgConfig())
		},
	}
	peerCmd = &cobra.Command{
		Use:   "peer",
		Short: "Peer Configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			for {
				selection = frontend.PeerMenuForm()
				switch selection {
				case frontend.Add:
					if err := peerAddCmd.RunE(cmd, args); err != nil {
						return err
					}
				case frontend.Edit:
					if err := peerEditCmd.RunE(cmd, args); err != nil {
						return err
					}
				case frontend.Delete:
					if err := peerDeleteCmd.RunE(cmd, args); err != nil {
						return err
					}
				case frontend.Quit:
					return nil
				}
			}
		},
	}

	peerAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add Peer",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _ = wgctrl.New()
			clientDevice, _ = client.Device(wirektun.DefaultTunName())
			if err := frontend.ValidateDeviceExists(client, wirektun.DefaultTunName()); err != nil {
				return err
			}
			config := frontend.PeerForm()
			return client.ConfigureDevice(clientDevice.Name, config)
		},
	}

	peerEditCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit Peer",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _ = wgctrl.New()
			clientDevice, _ = client.Device(wirektun.DefaultTunName())
			if err := frontend.ValidateDeviceExists(client, wirektun.DefaultTunName()); err != nil {
				return err
			}
			peer := frontend.PeerEditForm(clientDevice)
			return client.ConfigureDevice(clientDevice.Name, peer)
		},
	}

	peerDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete Peer",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _ = wgctrl.New()
			clientDevice, _ = client.Device(wirektun.DefaultTunName())
			if err := frontend.ValidateDeviceExists(client, wirektun.DefaultTunName()); err != nil {
				return err
			}
			peer := frontend.PeerDeleteForm(clientDevice)
			if peer == "" {
				return nil
			}
			config := wirekcfg.ToDeletePeerConfig(peer)
			return client.ConfigureDevice(clientDevice.Name, wgtypes.Config{
				Peers: []wgtypes.PeerConfig{config},
			})
		},
	}
)
