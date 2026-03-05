package cmd

import (
	"fmt"
	"time"
	"wirekcp/frontend"
	"wirekcp/wirekcfg"
	"wirekcp/wirektun"

	"github.com/inancgumus/screen"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/wirekcp/wgctrl"
)

var (
	keyCmd = &cobra.Command{
		Use:   "key",
		Short: "Show WireKCP key",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			client, err = wgctrl.New()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			clientDevice, err = client.Device(wirektun.DefaultTunName())
			if err != nil {
				return fmt.Errorf("failed to get device: %w", err)
			}

			if err := frontend.ValidateDeviceExists(client, wirektun.DefaultTunName()); err != nil {
				return err
			}

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
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !isAdmin() {
				cmd.SetHelpFunc(func(*cobra.Command, []string) {
					cmd.PrintErrln("Permission denied")
				})
				screen.Clear()
				screen.MoveTopLeft()
				return requestAdminCmd.Execute()
			}
			return performModeSwitch(cmd)
		},
	}
)

// performModeSwitch handles the actual mode switching logic
func performModeSwitch(cmd *cobra.Command) error {
	svc, err := newSVC(currentProgram, newSVCConfig())
	if err == service.ErrNotInstalled {
		cmd.PrintErrln("Service is not installed. Please run 'wirekcp service install' first.")
		return nil
	}
	if err != nil {
		return err
	}

	status, err := svc.Status()
	if err != nil {
		return err
	}

	if status == service.StatusRunning {
		if err := svc.Stop(); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}
		time.Sleep(2 * time.Second) // Wait for the service to stop completely
	}

	if status != service.StatusStopped {
		return fmt.Errorf("service is not in stopped state")
	}

	if err := switchModeInConfig(cmd); err != nil {
		return err
	}

	if err := svc.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	screen.Clear()
	screen.MoveTopLeft()
	cmd.Println(modeStyle.Render("Service started successfully!\n"))
	cmd.Println(modeStyle.Render("Press Enter to continue..."))
	fmt.Scanln()

	return nil
}

// switchModeInConfig toggles the mode in the config file
func switchModeInConfig(cmd *cobra.Command) error {
	config, err := wirekcfg.ReadFromFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	switch config.Mode {
	case "kcp":
		config.Mode = "udp"
	case "udp":
		config.Mode = "kcp"
	}

	if err := config.WriteToFile(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	screen.Clear()
	screen.MoveTopLeft()
	cmd.Println(successStyle.Render("WireKCP mode switched successfully!"))
	cmd.Println(modeStyle.Render(fmt.Sprintf("Switched to %s mode", config.Mode)))
	cmd.Println(modeStyle.Render("Starting service...\n"))

	return nil
}
