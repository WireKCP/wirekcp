package cmd

import (
	"fmt"
	"wirekcp/frontend"
	"wirekcp/wirekcfg"
	"wirekcp/wirektun"

	"github.com/charmbracelet/lipgloss"
	"github.com/inancgumus/screen"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/wirekcp/wgctrl"
)

var (
	keyCmd = &cobra.Command{
		Use:   "key",
		Short: "Show WireKCP key",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _ = wgctrl.New()
			clientDevice, _ = client.Device(wirektun.DefaultTunName())
			if err := frontend.ValidateDeviceExists(client, wirektun.DefaultTunName()); err != nil {
				return err
			}
			titleStyle := lipgloss.NewStyle().
				Width(52).
				Foreground(lipgloss.Color("#00FF00")).
				Align(lipgloss.Center).
				Padding(1, 1)
			contentStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Foreground(lipgloss.Color("#FFFFFF")).
				Width(50).
				Align(lipgloss.Center).
				Padding(0, 1)

			cmd.Println(titleStyle.Render("WireKCP Key"))
			cmd.Println(contentStyle.Render("Public Key"))
			cmd.Println(contentStyle.BorderTop(false).Render(clientDevice.PublicKey.String()))
			cmd.Println(contentStyle.Render("Private Key"))
			cmd.Println(contentStyle.BorderTop(false).Render(clientDevice.PrivateKey.String()))
			return nil
		},
	}
)

var (
	switchModeCmd = &cobra.Command{
		Use:   "switch",
		Short: "Switch WireKCP mode between KCP and UDP",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isAdmin() {
				cmd.SetHelpFunc(func(*cobra.Command, []string) {
					cmd.PrintErrln("Permission denied")
				})
				screen.Clear()
				screen.MoveTopLeft()
				requestAdminCmd.Execute()
				return nil
			}
			if serviceController, err := newSVC(&program{}, newSVCConfig()); err == service.ErrNotInstalled {
				cmd.PrintErrln("Service is not installed. Please run 'wirekcp service install' first.")
			} else if err != nil {
				return err
			} else {
				status, err := serviceController.Status()
				if err != nil {
					return err
				}
				if status == service.StatusRunning {
					cmd.PrintErrln("Please stop the service before switching modes.")
					return nil
				}
				if status == service.StatusStopped {
					var ifconfig *wirekcfg.Config
					ifconfig, err = wirekcfg.ReadFromFile(configPath)
					if err != nil {
						return fmt.Errorf("failed to read config file: %w", err)
					}
					if ifconfig.Mode == "kcp" {
						ifconfig.Mode = "udp"
					} else {
						ifconfig.Mode = "kcp"
					}
					if err = ifconfig.WriteToFile(configPath); err != nil {
						return fmt.Errorf("failed to write config file: %w", err)
					}

					successStyle := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#00FF00")).
						Background(lipgloss.Color("#222222")).
						Bold(true).
						Padding(1, 4).
						MarginTop(1).
						MarginBottom(1).
						Border(lipgloss.RoundedBorder()).
						BorderForeground(lipgloss.Color("#00FF00"))

					modeStyle := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#FFD700")).
						Bold(true)

					switchedTo := modeStyle.Render(fmt.Sprintf("Switched to %s mode", ifconfig.Mode))
					pleaseStartService := modeStyle.Render("Please start the service to apply changes.")

					cmd.Println(successStyle.Render("WireKCP mode switched successfully!"))
					cmd.Println(switchedTo)
					cmd.Println(pleaseStartService)
				}
			}
			return nil
		},
	}
)
