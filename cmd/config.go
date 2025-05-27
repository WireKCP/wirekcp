package cmd

import (
	"fmt"
	"runtime"
	"wirekcp/frontend"

	"github.com/charmbracelet/lipgloss"
	"github.com/inancgumus/screen"
	"github.com/spf13/cobra"
)

var (
	selection frontend.Selection
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Configuration Menu",
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
			for {
				switch frontend.MenuForm() {
				case frontend.Interface:
					if err := interfaceCmd.RunE(cmd, args); err != nil {
						return err
					}
				case frontend.Peer:
					if err := peerCmd.RunE(cmd, args); err != nil {
						return err
					}
				case frontend.SwitchMode:
					if err := switchModeCmd.RunE(cmd, args); err != nil {
						return err
					}
				case frontend.Quit:
					return nil
				}
			}
		},
	}
	requestAdminCmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			var admin string
			switch runtime.GOOS {
			case "windows":
				admin = "administrator"
			default:
				admin = "root"
			}
			titleStyle := lipgloss.NewStyle().
				Width(52).
				Foreground(lipgloss.Color("#00FF00")).
				Align(lipgloss.Center).
				Padding(1, 1)
			contentStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Foreground(lipgloss.Color("#FF0000")).
				Width(50).
				Align(lipgloss.Center).
				Padding(0, 1)

			cmd.PrintErrln(titleStyle.Render("Welcome to WireKCP !"))
			cmd.PrintErrln(contentStyle.Render(fmt.Sprintf("Please run as %s", admin)))
		},
	}
)
