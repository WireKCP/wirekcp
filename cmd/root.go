package cmd

import (
	"os"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wirekcp/wgctrl"
	"github.com/wirekcp/wgctrl/wgtypes"
)

const (
	// ExitSetupFailed defines exit code
	ExitSetupFailed   = 1
	DefaultConfigPath = ""
)

var (
	configPath        string
	defaultConfigPath string
	logLevel          string
	defaultLogFile    string
	logFile           string

	client       *wgctrl.Client
	clientDevice *wgtypes.Device

	rootCmd = &cobra.Command{
		Use: "wirekcp",
	}

	// Execution control channel for stopCh signal
	stopCh chan int
	term   chan os.Signal
	errs   chan error
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func defaultPort() uint16 {
	if s := os.Getenv("PORT"); s != "" {
		if p, err := strconv.ParseUint(s, 10, 16); err == nil {
			return uint16(p)
		}
	}
	if runtime.GOOS == "windows" {
		return 49200
	}
	return 0
}

func init() {
	stopCh = make(chan int)
	term = make(chan os.Signal, 1)
	errs = make(chan error)

	defaultConfigPath = "/etc/wirekcp/config.toml"
	defaultLogFile = "/var/log/wirekcp/wirekcp.log"
	if runtime.GOOS == "windows" {
		defaultConfigPath = os.Getenv("PROGRAMDATA") + "\\WireKCP\\" + "config.toml"
		defaultLogFile = os.Getenv("PROGRAMDATA") + "\\WireKCP\\" + "wirekcp.log"
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", defaultConfigPath, "WireKCP config file location")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "sets WireKCP log level")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", defaultLogFile, "sets WireKCP log path. If console is specified the the log will be output to stdout")
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(interfaceCmd, peerCmd)
	peerCmd.AddCommand(peerAddCmd, peerEditCmd, peerDeleteCmd)
	serviceCmd.AddCommand(runCmd, startCmd, stopCmd, restartCmd) // service control commands are subcommands of service
	serviceCmd.AddCommand(installCmd, uninstallCmd)              // service installer commands are subcommands of service
}
