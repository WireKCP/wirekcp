package cmd

import (
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	installCmd = &cobra.Command{
		Use:   "install",
		Short: "installs wirekcp service",
		Run: func(cmd *cobra.Command, args []string) {

			svcConfig := newSVCConfig()

			svcConfig.Arguments = []string{
				"service",
				"run",
				"--config",
				configPath,
				"--log-level",
				logLevel,
			}

			if runtime.GOOS == "linux" {
				// Respected only by systemd systems
				svcConfig.Dependencies = []string{"After=network.target syslog.target"}
			} else if runtime.GOOS == "windows" && !isAdmin() {
				runAsAdmin()
				args := strings.Join(os.Args[1:], " ")
				cmd.Printf("Starting following command as administrator: %s %s\n", os.Args[0], args)
				return
			}

			s, err := newSVC(&program{}, svcConfig)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			err = s.Install()
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Printf("Wirekcp service has been installed")
		},
	}
)

var (
	uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "uninstalls wirekcp service from system",
		Run: func(cmd *cobra.Command, args []string) {
			if runtime.GOOS == "windows" && !isAdmin() {
				runAsAdmin()
				return
			}

			s, err := newSVC(&program{}, newSVCConfig())
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			err = s.Uninstall()
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Printf("Wirekcp has been uninstalled")
		},
	}
)

func init() {
}
