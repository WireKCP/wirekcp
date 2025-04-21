package cmd

import (
	"fmt"

	"wirekcp/wgengine"
	"wirekcp/wirektun"

	"github.com/spf13/cobra"
	"github.com/wirekcp/wireguard-go/device"
)

var (
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Start WireKCP",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			interfaceName := wirektun.DefaultTunName()
			logger := device.NewLogger(
				logIntLevel,
				fmt.Sprintf("(%s) ", interfaceName),
			)
			var engine wgengine.Engine
			engine, err = wgengine.NewUserspaceEngine(logger, interfaceName, defaultPort(), configPath, logFile)
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
