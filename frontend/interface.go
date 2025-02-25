package frontend

import (
	"bytes"
	"strconv"
	"wirekcp/wirekcfg"
	"wirekcp/wirektun"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/wirekcp/wgctrl/wgtypes"
)

func InterfaceForm(device wgtypes.Device) *wirekcfg.Config {
	var config wirekcfg.Config
	config.IPv4CIDR = wirektun.GetIP()
	config.ListenPort = device.ListenPort
	emptyKey := wgtypes.Key{}
	if !bytes.Equal(device.PrivateKey[:], emptyKey[:]) {
		config.PrivateKey = device.PrivateKey.String()
	}
	port := strconv.Itoa(config.ListenPort)
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("IP Address <CIDR>").Validate(ValidateCIDR).Value(&config.IPv4CIDR),
			huh.NewInput().Title("Listen Port").Validate(ValidateRequiredInt).Value(&port),
			huh.NewInput().Title("Private Key").Validate(ValidateOptionalKey).Value(&config.PrivateKey),
		),
	).WithProgramOptions(tea.WithAltScreen()).Run(); err != nil {
		return nil
	}
	uport, _ := strconv.Atoi(port)
	config.ListenPort = uport
	if config.PrivateKey == "" {
		key, _ := wgtypes.GeneratePrivateKey()
		config.PrivateKey = key.String()
	}
	return &config
}
