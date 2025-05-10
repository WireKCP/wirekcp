package cmd

import (
	"fmt"

	"wirekcp/wgengine"
	"wirekcp/wireklog"
	"wirekcp/wirektun"
	"wirekcp/wirekutils"

	"github.com/spf13/cobra"
	"github.com/wirekcp/wireguard-go/device"
)

var (
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Start WireKCP",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var logger *device.Logger
			interfaceName := wirektun.DefaultTunName()
			if foreground {
				logger = wireklog.NewStdoutLogger(
					logIntLevel,
					fmt.Sprintf("(%s) ", interfaceName),
				)
			} else {
				fd, err := wireklog.OpenOrCreateFile(logFile)
				if err != nil {
					return err
				}
				defer fd.Close()
				logger = wireklog.NewFileLogger(
					logIntLevel,
					fmt.Sprintf("(%s) ", interfaceName),
					fd,
				)
			}
			if err := wirekutils.CheckOrInstallWinTun(); err != nil {
				logger.Errorf("Failed to check or install Wintun: %v", err)
				return err
			}
			var engine wgengine.Engine
			engine, err = wgengine.NewUserspaceEngine(logger, interfaceName, defaultPort(), configPath)
			if err != nil {
				logger.Errorf("Failed to create userspace engine: %v", err)
				return err
			}
			defer engine.Close()

			setupCloseHandler(logger)

			select {
			case <-stopCh:
			case <-errs:
			case <-engine.Wait():
			}

			// clean up
			engine.Close()
			return nil
		},
	}
)
