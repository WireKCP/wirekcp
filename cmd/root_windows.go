//go:build windows

package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wirekcp/wireguard-go/device"
	"golang.org/x/sys/windows"
)

func setupCloseHandler(logger *device.Logger) {
	signal.Notify(term, os.Interrupt, windows.SIGTERM, windows.SIGINT, windows.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for range term {
			logger.Verbosef("Received SIGTERM signal")
			logger.Verbosef("Shutting down WireKCP")
			logger.Verbosef("Ctrl+C pressed in Terminal")
			stopCh <- 0
		}
	}()
}
