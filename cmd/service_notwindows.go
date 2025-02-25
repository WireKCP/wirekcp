//go:build !windows

package cmd

import "os"

func isAdmin() bool {
	if os.Getuid() == 0 {
		return true
	}
	return false
}

func runAsAdmin() error {
	return nil
}
