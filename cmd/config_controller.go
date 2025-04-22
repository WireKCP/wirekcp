package cmd

import (
	"errors"
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
				return errors.New("invalid configuration")
			}
			if err := client.ConfigureDevice(clientDevice.Name, *config.ToWgConfig()); err != nil {
				return err
			}
			if err := wirekcfg.SetIPwithoutTun(config.IPv4CIDR); err != nil {
				return err
			}
			return config.ChangeInterface(configPath)
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
			if err := client.ConfigureDevice(clientDevice.Name, config); err != nil {
				return err
			}
			fileConfig, err := wirekcfg.ReadFromFile(configPath)
			if err != nil {
				return err
			}
			fileConfig.Peers = wirekcfg.ToWkPeersConfig(config.Peers)
			return fileConfig.WriteToFile(configPath)
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
			config := frontend.PeerEditForm(clientDevice)
			if err := client.ConfigureDevice(clientDevice.Name, config); err != nil {
				return err
			}
			fileConfig, err := wirekcfg.ReadFromFile(configPath)
			if err != nil {
				return err
			}
			fileConfig.Peers = wirekcfg.ToWkPeersConfig(config.Peers)
			return fileConfig.WriteToFile(configPath)
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
			if err := client.ConfigureDevice(clientDevice.Name, wgtypes.Config{
				Peers: []wgtypes.PeerConfig{config},
			}); err != nil {
				return err
			}
			fileConfig, err := wirekcfg.ReadFromFile(configPath)
			if err != nil {
				return err
			}
			for i, p := range fileConfig.Peers {
				if p.PublicKey == peer {
					fileConfig.Peers = append(fileConfig.Peers[:i], fileConfig.Peers[i+1:]...)
					break
				}
			}
			return fileConfig.WriteToFile(configPath)
		},
	}
)
