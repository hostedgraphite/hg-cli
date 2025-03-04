package tui

import (
	"fmt"
	"os"

	"github.com/hostedgraphite/hg-cli/sysinfo"
	"github.com/hostedgraphite/hg-cli/tui/types"
	"github.com/hostedgraphite/hg-cli/tui/views"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentView types.View
	sysinfo     sysinfo.SysInfo
}

func (m model) Init() tea.Cmd {
	return m.currentView.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.currentView, cmd = m.currentView.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.currentView.View()
}

func StartTui(sysinfo sysinfo.SysInfo) error {

	p := tea.NewProgram(model{currentView: &views.MainView{SysInfo: sysinfo}, sysinfo: sysinfo})
	_, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		return err
	}

	return err
}
