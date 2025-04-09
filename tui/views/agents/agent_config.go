package agents

import (
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf"

	"github.com/hostedgraphite/hg-cli/styles"
	"github.com/hostedgraphite/hg-cli/sysinfo"
	"github.com/hostedgraphite/hg-cli/tui/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type AgentConfigView struct {
	form            *huh.Form
	agent           string
	action          string
	apiKey          string
	selectedPlugins []string
	sysInfo         sysinfo.SysInfo
	serviceSettings map[string]string
}

func NewAgentConfigView(agent, action string, sysInfo sysinfo.SysInfo) *AgentConfigView {
	var err error
	var actionGroup *huh.Group
	var apikey, selectedInstall, path string
	var selectedPlugins []string
	var confirmUninstall bool
	settings := telegraf.GetServiceSettings(sysInfo.Os, sysInfo.Arch, sysInfo.PkgMngr)

	agentViews := NewAgentsFields(agent)

	switch action {
	case "Install":
		actionGroup, err = agentViews.InstallView(apikey, selectedInstall, selectedPlugins)
	case "Uninstall":
		actionGroup, err = agentViews.UninstallView(confirmUninstall)
	case "Update Api Key":
		actionGroup, err = agentViews.UpdateApiKeyView(apikey, path, settings["configPath"])
	default:
		return nil
	}

	if err != nil {
		return nil
	}

	form := huh.NewForm(actionGroup).
		WithWidth(80).
		WithTheme(styles.AgentsPageStyle(agent)).
		WithHeight(30).
		WithKeyMap(styles.CustomKeyMap())

	return &AgentConfigView{
		form:            form,
		agent:           agent,
		action:          action,
		apiKey:          apikey,
		sysInfo:         sysInfo,
		serviceSettings: settings,
	}
}
func (a *AgentConfigView) Init() tea.Cmd {

	if a.form == nil {
		return nil
	}
	return tea.Batch(a.form.Init(), tea.ClearScreen, tea.EnterAltScreen)
}

func (a *AgentConfigView) Update(msg tea.Msg) (types.View, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.sysInfo.Width = msg.Width
		a.sysInfo.Height = msg.Height
		return a, tea.Batch(tea.ClearScreen, tea.EnterAltScreen)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "esc", "q":
			return a, tea.Quit
		case "enter":
		}
	}

	if a.form != nil {
		form, cmd := a.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			a.form = f
			cmds = append(cmds, cmd)
		}
	}

	if a.form.State == huh.StateCompleted {

		options := map[string]interface{}{}

		a.apiKey = a.form.GetString("apikey")
		options["apikey"] = a.apiKey
		switch a.action {
		case "Install":
			if a.agent == "Telegraf" {
				plugins := a.form.Get("plugins")
				if val, ok := plugins.([]string); ok && len(val) > 0 {
					a.selectedPlugins = val
				}
				options["plugins"] = a.selectedPlugins
			}

		case "Update Api Key":
			path := a.form.GetString("path")
			if path == "" {
				path = a.serviceSettings["configPath"]
			}
			options["config"] = path
		case "Uninstall":
			if !a.form.GetBool("confirmUninstall") {
				return a, tea.Quit
			}
		}

		agentRunner := NewAgentRunner(a.agent, a.action, options, a.sysInfo, a.serviceSettings)
		return agentRunner, agentRunner.Init()
	}

	return a, tea.Batch(cmds...)

}
func (a *AgentConfigView) View() string {
	if a.form == nil {
		return "Form not intialized"
	}
	switch a.form.State {
	case huh.StateCompleted:
		return "completed"
	default:
		return styles.PlaceContent(
			a.sysInfo.Width,
			a.sysInfo.Height,
			(styles.DefaultStyles().Page.Render(a.form.View())))
	}
}
