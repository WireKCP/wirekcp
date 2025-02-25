package cmd

import (
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type program struct {
	cmd  *cobra.Command
	args []string
}

func newSVCConfig() *service.Config {
	return &service.Config{
		Name:        "wirekcp",
		DisplayName: "WireKCP",
		Description: "A KCP-based WireGuard network that connects your devices into a single private network.",
	}
}

func newSVC(prg *program, conf *service.Config) (service.Service, error) {
	s, err := service.New(prg, conf)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return s, nil
}

var (
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "manages WireKCP service",
	}
)

func init() {
}
