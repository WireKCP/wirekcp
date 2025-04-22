package wireklog

import (
	"log"
	"os"
	"path/filepath"

	chramLog "github.com/charmbracelet/log"
	"github.com/wirekcp/wireguard-go/device"
)

func OpenOrCreateFile(path string) (*os.File, error) {
	var fd *os.File
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If the Directory does not exist, create it
		folderPath := filepath.Dir(path)
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			return nil, err
		}
		// Create the file if it does not exist
		fd, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	} else {
		fd, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}
	return fd, nil
}

func NewFileLogger(level int, prepend string, fd *os.File) *device.Logger {
	logger := &device.Logger{Verbosef: device.DiscardLogf, Errorf: device.DiscardLogf}
	logf := func(prefix string) func(string, ...any) {
		return log.New(fd, prefix+": "+prepend, log.Ldate|log.Ltime).Printf
	}
	if level >= device.LogLevelVerbose {
		logger.Verbosef = logf("DEBUG")
	}
	if level >= device.LogLevelError {
		logger.Errorf = logf("ERROR")
	}
	return logger
}

func NewStdoutLogger(level int, prepend string) *device.Logger {
	logger := &device.Logger{Verbosef: device.DiscardLogf, Errorf: device.DiscardLogf}
	stdLogger := chramLog.NewWithOptions(os.Stdout, chramLog.Options{Prefix: prepend, ReportTimestamp: true, Level: chramLog.DebugLevel})
	if level >= device.LogLevelVerbose {
		logger.Verbosef = stdLogger.Debugf
	}
	if level >= device.LogLevelError {
		logger.Errorf = stdLogger.Errorf
	}
	return logger
}
