//go:build windows

package wirekutils

import (
	"os"
)

func CheckOrInstallWinTun() error {
	filename := "wintun.dll"
	destination := "C:\\Windows\\System32\\wintun.dll"
	// Check if the file exists, if not link it
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		cwd, err := os.Getwd()
		println("Current working directory: " + cwd)
		if err != nil {
			return err
		}
		file := cwd + "\\" + filename
		println("Linking " + file + " to " + destination)
		// Link the file to the destination
		return os.Link(file, destination)
	}
	return nil
}
