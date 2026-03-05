package cmd

import (
	"runtime"
	"time"

	"github.com/charmbracelet/log"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	log.Infof("Starting service %s", s.String())
	go p.run()
	return nil
}

func (p *program) run() {
	// Set current program reference for upCmd to access
	currentProgram = p
	defer func() { currentProgram = nil }()

	// Wait briefly on Windows to allow network drivers to fully initialize
	// This helps prevent "device not ready" errors on startup
	if runtime.GOOS == "windows" {
		log.Infof("Waiting for network drivers to initialize...")
		time.Sleep(5 * time.Second)
	}

	// Run the service
	err := upCmd.RunE(p.cmd, p.args)
	if err != nil {
		s, _ := newSVC(p, newSVCConfig())
		p.Stop(s)
		return
	}
}

func (p *program) Stop(s service.Service) error {
	log.Infof("Stopping service %s", s.String())
	// Stop the engine gracefully using context cancellation
	if p.engine != nil {
		p.engine.Close()
	}
	return nil
}

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "manually runs WireKCP in service mode (for debugging)",
		Run: func(cmd *cobra.Command, args []string) {
			prg := &program{
				cmd:  cmd,
				args: args,
			}
			s, err := newSVC(prg, newSVCConfig())
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			err = s.Run()
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		},
	}
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "starts wirekcp service",
		Run: func(cmd *cobra.Command, args []string) {
			if runtime.GOOS == "windows" && !isAdmin() {
				runAsAdmin()
				cmd.Println("Wirekcp service has been started")
				return
			}

			s, err := newSVC(&program{}, newSVCConfig())
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			err = s.Start()
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Println("Wirekcp service has been started")
		},
	}
)

var (
	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "stops wirekcp service",
		Run: func(cmd *cobra.Command, args []string) {
			if runtime.GOOS == "windows" && !isAdmin() {
				runAsAdmin()
				cmd.Println("Wirekcp service has been stopped")
				return
			}

			s, err := newSVC(&program{}, newSVCConfig())
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			err = s.Stop()
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Println("Wirekcp service has been stopped")
		},
	}
)

var (
	restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "restarts wirekcp service",
		Run: func(cmd *cobra.Command, args []string) {

			s, err := newSVC(&program{}, newSVCConfig())
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			err = s.Restart()
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Println("Wirekcp service has been restarted")
		},
	}
)

func init() {
}
