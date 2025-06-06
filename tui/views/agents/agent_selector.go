package agents

import (
	"fmt"
	"slices"

	"github.com/hostedgraphite/hg-cli/styles"
	"github.com/hostedgraphite/hg-cli/sysinfo"
	"github.com/hostedgraphite/hg-cli/tui/types"
	"github.com/hostedgraphite/hg-cli/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

var agents = []string{"Telegraf", "OpenTelemetry"}
var commingSoon = []string{""}
var agentActions = []string{"Install", "Update Api Key", "Uninstall"}

type AgentsView struct {
	sysInfo sysinfo.SysInfo
	form    *huh.Form
	agent   string
	action  string
}

func NewAgentView(sysInfo sysinfo.SysInfo) *AgentsView {
	var selectedAgent, selectedAction string

	// First form group: Select an agent
	agentActionGroup := huh.NewGroup(
		huh.NewNote().
			Title(styles.AgentsPageTitle),

		// agent selection
		huh.NewSelect[string]().
			Key("agent").
			Title("Select Agent").
			Description("Select any of the following agents, and we will guide you through their installation").
			Options(huh.NewOptions(agents...)...).
			Validate(func(agent string) error {
				if slices.Contains(commingSoon, selectedAgent) {
					return fmt.Errorf("sorry, this agent is not yet available")
				}
				return nil
			}).
			Value(&selectedAgent),

		// agent action
		huh.NewSelect[string]().
			Key("action").
			Title("Select Agent Action").
			Description("Choose one of the following actions: Install, Update your Api key, Uninstall").
			OptionsFunc(func() []huh.Option[string] {
				if slices.Contains(commingSoon, selectedAgent) {
					return huh.NewOptions([]string{"Comming Soon"}...)
				} else {
					return huh.NewOptions(agentActions...)
				}
			}, &selectedAgent).
			Validate(func(action string) error {
				if utils.AgentRequiresSudo(sysInfo.Os, action, sysInfo.PkgMngr, selectedAgent) && !sysInfo.SudoPerm {
					return fmt.Errorf("this action requires admin privileges, please run as root")
				}
				return nil
			}).
			Value(&selectedAction),
	)

	form := huh.NewForm(agentActionGroup).
		WithTheme(styles.FormStyles()).
		WithWidth(80).
		WithHeight(30).
		WithKeyMap(styles.CustomKeyMap())

	return &AgentsView{
		form:    form,
		agent:   selectedAgent,
		action:  selectedAction,
		sysInfo: sysInfo,
	}
}

func (a *AgentsView) Init() tea.Cmd {
	if a.form == nil {
		return nil
	}

	return tea.Batch(a.form.Init(), tea.ClearScreen, tea.EnterAltScreen)
}

func (a *AgentsView) Update(msg tea.Msg) (types.View, tea.Cmd) {
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
		// The last selected value of the form doesn't update correctly
		// so we use a Key() value in the form to then get the value.
		a.agent = a.form.GetString("agent")
		a.action = a.form.GetString("action")
		configurationView := NewAgentConfigView(a.agent, a.action, a.sysInfo)

		return configurationView, configurationView.Init()
	}

	return a, tea.Batch(cmds...)
}

func (a *AgentsView) View() string {
	if a.form == nil {
		return "Form not initialized"
	}

	switch a.form.State {
	case huh.StateCompleted:
		return styles.PlaceContent(
			a.sysInfo.Width,
			a.sysInfo.Height,
			styles.DefaultStyles().Page.Render(
				(fmt.Sprintf("Agent: %s, Action: %s", a.agent, a.action))))
	default:
		return styles.PlaceContent(
			a.sysInfo.Width,
			a.sysInfo.Height,
			(styles.DefaultStyles().Page.Render(a.form.View())))
	}
}
