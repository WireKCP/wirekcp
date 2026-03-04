package cmd

import (
	"fmt"

	"wirekcp/wgengine"
	"wirekcp/wireklog"
	"wirekcp/wirektun"

	"github.com/spf13/cobra"
	"github.com/wirekcp/wireguard-go/device"
)

var (
	// currentProgram holds reference to current program instance for service mode
	currentProgram *program
	serviceMode    bool = false
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
			var engine wgengine.Engine
			engine, err = wgengine.NewUserspaceEngine(logger, interfaceName, defaultPort(), configPath)
			if err != nil {
				logger.Errorf("Failed to create userspace engine: %v", err)
				return err
			}
			defer engine.Close()

			// Store engine reference if running as a service
			if currentProgram != nil {
				currentProgram.engine = engine
			}

			setupCloseHandler(logger)

			select {
			case <-stopCh:
			case <-errs:
			case <-engine.Wait():
			}

			engine.Close()
			return nil
		},
	}
)
