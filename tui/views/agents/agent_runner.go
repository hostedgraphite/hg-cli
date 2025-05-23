package agents

import (
	"strings"

	"github.com/hostedgraphite/hg-cli/agentmanager"
	"github.com/hostedgraphite/hg-cli/formatters"
	"github.com/hostedgraphite/hg-cli/pipeline"
	"github.com/hostedgraphite/hg-cli/styles"
	"github.com/hostedgraphite/hg-cli/sysinfo"
	"github.com/hostedgraphite/hg-cli/tui/types"

	"github.com/charmbracelet/bubbles/spinner"
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
	runner          *pipeline.Runner
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

	switch a.action {
	case "Install":
		agent := agentmanager.NewAgent(a.agent, a.options, a.sysInfo)
		updates := make(chan *pipeline.Pipe)
		installPipeline, err := agent.InstallPipeline(updates)
		if err != nil {
			panic(err) // This BAD. TODO: not this
		}
		runner := pipeline.NewRunner(
			installPipeline,
			false,
			updates,
		)
		a.runner = runner
		return a.runner.RunStatic()
	case "Update Api Key":
		agent := agentmanager.NewAgent(a.agent, a.options, a.sysInfo)
		updates := make(chan *pipeline.Pipe)
		updateApikeyPipeline, err := agent.UpdateApiKeyPipeline(updates)
		if err != nil {
			panic(err) // This BAD. TODO: not this
		}
		runner := pipeline.NewRunner(
			updateApikeyPipeline,
			false,
			updates,
		)
		a.runner = runner
		return a.runner.RunStatic()
	case "Uninstall":
		agent := agentmanager.NewAgent(a.agent, nil, a.sysInfo)
		updates := make(chan *pipeline.Pipe)
		uninstallPipeline, err := agent.UninstallPipeline(updates)
		if err != nil {
			panic(err) // This BAD. TODO: not this
		}
		runner := pipeline.NewRunner(
			uninstallPipeline,
			false,
			updates,
		)
		a.runner = runner
		return a.runner.RunStatic()
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
	case pipeline.PipeUpdate:
		var cmd tea.Cmd
		_, cmd = a.runner.Update(msg)
		cmds = append(cmds, cmd)
	case spinner.TickMsg:
		var cmd tea.Cmd
		_, cmd = a.runner.Update(msg)
		cmds = append(cmds, cmd)
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
	var summary formatters.SummaryContent
	if a.runner != nil {
		if a.runner.Pipeline.IsCompleted() {
			switch a.action {
			case "Install":
				if a.agent == "Telegraf" {
					summary = &formatters.TelegrafSummary{
						ActionSummary: formatters.ActionSummary{
							Agent:    a.agent,
							Success:  a.runner.Pipeline.Success(),
							Action:   a.action,
							Config:   a.serviceSettings["configPath"],
							StartCmd: a.serviceSettings["startHint"],
						},
						Plugins: a.options["plugins"].([]string),
					}
				} else if a.agent == "OpenTelemetry" {
					summary = &formatters.OtelContribSummary{
						ActionSummary: formatters.ActionSummary{
							Agent:    a.agent,
							Success:  a.runner.Pipeline.Success(),
							Action:   a.action,
							Config:   a.serviceSettings["configPath"],
							StartCmd: a.serviceSettings["startHint"],
						},
						Receiver: "hostmetrics",
						Exporter: "carbon",
					}
				}
			case "Update Api Key":
				data := formatters.ActionSummary{
					Agent:      a.agent,
					Success:    a.runner.Pipeline.Success(),
					Action:     a.action,
					Config:     a.options["config"].(string),
					RestartCmd: a.serviceSettings["restartHint"],
				}

				if a.agent == "Telegraf" {
					summary = &formatters.TelegrafSummary{
						ActionSummary: data,
					}
				} else if a.agent == "OpenTelemetry" {
					summary = &formatters.OtelContribSummary{
						ActionSummary: data,
					}
				}
			case "Uninstall":
				data := formatters.ActionSummary{
					Agent:   a.agent,
					Success: a.runner.Pipeline.Success(),
					Action:  a.action,
				}

				if a.agent == "Telegraf" {
					summary = &formatters.TelegrafSummary{
						ActionSummary: data,
					}
				} else if a.agent == "OpenTelemetry" {
					summary = &formatters.OtelContribSummary{
						ActionSummary: data,
					}
				}
			}
		} else {
			s := styles.DefaultStyles()
			content := a.runner.View()
			content = s.Runner.Render(content)
			return styles.PlaceContent(
				a.sysInfo.Width,
				a.sysInfo.Height,
				styles.DefaultStyles().Page.Render(content),
			)
		}
	}

	return formatters.GenerateSummary(summary, a.sysInfo.Width, a.sysInfo.Height)
}
