//go:build !windows

package wirekutils

func CheckOrInstallWinTun() error {
	return nil
}
