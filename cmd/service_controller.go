package cmd

import (
	"runtime"

	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	log.Info("starting service") //nolint
	go p.run()
	return nil
}

func (p *program) run() {
	// Run the service
	err := upCmd.RunE(p.cmd, p.args)
	if err != nil {
		stopCh <- 0
		return
	}
}

func (p *program) Stop(s service.Service) error {
	stopCh <- 1
	return nil
}

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "runs wirekcp as service",
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
			cmd.Printf("Wirekcp service is running")
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
			cmd.Printf("Wirekcp service has been started")
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
			cmd.Printf("Wirekcp service has been stopped")
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
			cmd.Printf("Wirekcp service has been restarted")
		},
	}
)

func init() {
}
