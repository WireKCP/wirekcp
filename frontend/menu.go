package frontend

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func MenuForm(mode string) Selection {
	var selection Selection
	huh.NewForm(
		huh.NewGroup(huh.NewSelect[Selection]().Title("Main Menu").Options(
			huh.NewOption("Interface Configuration", Interface),
			huh.NewOption("Peer Configuration", Peer),
			huh.NewOption(fmt.Sprintf("Switch Mode '%s'", mode), SwitchMode),
			huh.NewOption("Quit", Quit)).Value(&selection),
		)).WithProgramOptions(tea.WithAltScreen()).Run()
	return selection
}
