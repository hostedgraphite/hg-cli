package agents

import (
	"strings"
	"time"

	"github.com/hostedgraphite/hg-cli/agentmanager"
	"github.com/hostedgraphite/hg-cli/formatters"
	"github.com/hostedgraphite/hg-cli/styles"
	"github.com/hostedgraphite/hg-cli/sysinfo"
	"github.com/hostedgraphite/hg-cli/tui/types"

	tea "github.com/charmbracelet/bubbletea"
)

type updateMsg struct {
	current string
	rest    chan string
}

func (u updateMsg) awaitNext() updateMsg {
	return updateMsg{
		current: <-u.rest,
		rest:    u.rest,
	}
}

type AgentRunner struct {
	agent           string
	action          string
	options         map[string]interface{}
	sysInfo         sysinfo.SysInfo
	currentUpdate   string
	serviceSettings map[string]string
}

func NewAgentRunner(agent, action string, options map[string]interface{}, sysInfo sysinfo.SysInfo, serviceSettings map[string]string) *AgentRunner {
	return &AgentRunner{
		agent:           agent,
		action:          action,
		options:         options,
		sysInfo:         sysInfo,
		serviceSettings: serviceSettings,
	}
}

func (a *AgentRunner) Init() tea.Cmd {
	updates := make(chan string)
	agent := agentmanager.GetAgent(a.agent)

	apikey := a.options["apikey"].(string)
	switch a.action {
	case "Install":
		go func() {
			time.Sleep(2 * time.Second)
			agent.Install(apikey, a.sysInfo, a.options, updates)
		}()
	case "Update Api Key":
		go func() {
			time.Sleep(2 * time.Second)
			agent.UpdateApiKey(apikey, a.sysInfo, a.options, updates)
		}()
	case "Uninstall":
		go func() {
			time.Sleep(2 * time.Second)
			agent.Uninstall(a.sysInfo, updates)
		}()
	}

	return func() tea.Msg {
		return updateMsg{<-updates, updates}
	}
}

func (a *AgentRunner) Update(msg tea.Msg) (types.View, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.sysInfo.Width = msg.Width
		a.sysInfo.Width = msg.Height
		return a, tea.Batch(tea.ClearScreen, tea.EnterAltScreen)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "esc", "q":
			return a, tea.Quit
		case "enter":
		}
	case updateMsg:
		a.currentUpdate = msg.current
		if strings.HasPrefix(a.currentUpdate, "Error") {
			a.currentUpdate = "Error"
			return a, tea.Batch(cmds...)
		}

		if strings.HasPrefix(a.currentUpdate, "Completed") {
			return a, tea.Batch(cmds...)
		}

		return a, func() tea.Msg { return msg.awaitNext() }
	}

	return a, tea.Batch(cmds...)
}

func (a *AgentRunner) View() string {
	var summary formatters.ActionSummary
	var generatedSummary string

	if strings.HasPrefix(a.currentUpdate, "Completed") {
		if a.action == "Install" {
			summary = formatters.ActionSummary{
				Agent:    a.agent,
				Success:  true,
				Action:   a.action,
				Plugins:  a.options["plugins"].([]string),
				Config:   a.serviceSettings["configPath"],
				StartCmd: a.serviceSettings["startHint"],
			}
		} else if a.action == "Update Api Key" {
			summary = formatters.ActionSummary{
				Agent:      a.agent,
				Success:    true,
				Action:     a.action,
				Config:     a.options["config"].(string),
				RestartCmd: a.serviceSettings["restartHint"],
			}
		} else {
			// Uninstall
			summary = formatters.ActionSummary{
				Agent:   a.agent,
				Success: true,
				Action:  a.action,
			}
		}

		generatedSummary = formatters.GenerateSummary(summary, a.sysInfo.Width, a.sysInfo.Height)
	} else {
		s := styles.DefaultStyles()
		content := s.Updates.Render(a.currentUpdate)
		content = s.Page.Render(content)
		return styles.PlaceContent(
			a.sysInfo.Width,
			a.sysInfo.Height,
			content,
		)
	}

	return generatedSummary
}
