package cmd

import (
	"wirekcp/frontend"
	"wirekcp/wirektun"

	"github.com/charmbracelet/lipgloss"
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
