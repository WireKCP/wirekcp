//go:build !windows

package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wirekcp/wireguard-go/device"
)

func setupCloseHandler(logger *device.Logger) {
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range term {
			logger.Verbosef("Received SIGTERM signal")
			logger.Verbosef("Shutting down WireKCP")
			logger.Verbosef("Ctrl+C pressed in Terminal")
			stopCh <- 0
		}
	}()
}
