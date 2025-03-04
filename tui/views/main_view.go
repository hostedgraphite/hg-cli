package views

import (
	"strings"

	"github.com/hostedgraphite/hg-cli/styles"
	"github.com/hostedgraphite/hg-cli/sysinfo"
	"github.com/hostedgraphite/hg-cli/tui/types"
	"github.com/hostedgraphite/hg-cli/tui/views/agents"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var mainOptions = []string{"Agents", "Exit"}

type MainView struct {
	selectedIndex int
	SysInfo       sysinfo.SysInfo
}

func (m MainView) Init() tea.Cmd {
	return nil
}

func (m MainView) Update(msg tea.Msg) (types.View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SysInfo.Width = msg.Width
		m.SysInfo.Height = msg.Height
		return m, tea.Batch(tea.ClearScreen, tea.EnterAltScreen)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "l", "right":
			if m.selectedIndex < len(mainOptions)-1 {
				m.selectedIndex++
			}
		case "h", "left":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "enter":
			if m.selectedIndex == 0 {
				agentsView := agents.NewAgentView(m.SysInfo)
				return agentsView, agentsView.Init()
			} else if m.selectedIndex == 1 {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m MainView) View() string {
	s := styles.DefaultStyles()
	var styledOptions []string

	var viewStr strings.Builder
	description := "Welcome to the MetricFire CLI – your command-line sidekick for managing MetricFire like a pro! Right now, it’s all about making Agent Configuration Management a breeze, but stay tuned—more powerful tools and integrations are on the way! 🚀"
	viewStr.WriteString(s.Title.Render(styles.MetricfireLogo))
	viewStr.WriteString(s.MenuTitle.Render("Main Menu"))
	viewStr.WriteString(s.Description.Render(description))

	for i, option := range mainOptions {
		if i == m.selectedIndex {
			styledOptions = append(styledOptions, s.Selected.Render("> "+option))
		} else {
			styledOptions = append(styledOptions, s.Selections.Render("  "+option))
		}
	}

	menuRow := lipgloss.JoinHorizontal(lipgloss.Center, styledOptions...)
	viewStr.WriteString(menuRow)

	viewStr.WriteString(s.Footer.Render("\nPress ←/→ to move, Enter to select, q to quit."))
	content := s.Page.Render(viewStr.String())

	return styles.PlaceContent(m.SysInfo.Width, m.SysInfo.Height, content)
}
