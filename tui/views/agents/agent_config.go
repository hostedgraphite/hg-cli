package agents

import (
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf"

	"github.com/hostedgraphite/hg-cli/styles"
	"github.com/hostedgraphite/hg-cli/sysinfo"
	"github.com/hostedgraphite/hg-cli/tui/types"
	"github.com/hostedgraphite/hg-cli/tui/views/config"

	"github.com/hostedgraphite/hg-cli/utils"

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
	var actionGroup *huh.Group
	var apikey, selectedInstall, path string
	var selectedPlugins []string
	var confirmUninstall bool
	settings := telegraf.GetServiceSettings(sysInfo.Os, sysInfo.Arch, sysInfo.PkgMngr)

	header := getHeader(agent)
	if action == "Install" {
		actionGroup = installGroup(header, agent, apikey, selectedInstall, selectedPlugins)
	} else if action == "Update Api Key" {
		actionGroup = updateAPIKeyGroup(header, apikey, path, settings["configPath"])
	} else if action == "Uninstall" {
		actionGroup = uninstallGroup(header, confirmUninstall)
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
			// Install agent
			// Next few steps are in the cases where the options aren't binded to the stuct fields
			// This also sets the default values for the Telegraf default install option.
			installType := a.form.GetString("installType")
			if installType == "Default" && a.agent == "Telegraf" {
				options["plugins"] = telegraf.DefaultTelegrafPlugins
			} else {
				plugins := a.form.Get("plugins")
				if val, ok := plugins.([]string); ok && len(val) > 0 {
					a.selectedPlugins = val
				}
				options["plugins"] = a.selectedPlugins
			}

		case "Update Api Key":
			// Update agent
			path := a.form.GetString("path")
			if path == "" {
				path = a.serviceSettings["configPath"]
			}
			options["config"] = path
		case "Uninstall":
			if !a.form.GetBool("confirm") {
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

func getHeader(agent string) string {
	switch agent {
	case "Telegraf":
		return styles.MfAndTelegrafTitle
	default:
		return styles.MetricfireLogo
	}
}

func installGroup(header, agent, apikey, selectedInstall string, selectedPlugins []string) *huh.Group {

	actionGroup := huh.NewGroup(
		huh.NewNote().
			Title(header),

		huh.NewInput().
			Key("apikey").
			Title("Enter your Hosted Graphite API KEY").
			Prompt("API Key: ").
			Validate(func(s string) error {
				err := utils.ValidateAPIKey(apikey)
				if err != nil {
					return err
				}
				return nil
			}).
			Value(&apikey).
			EchoMode(huh.EchoModePassword),

		huh.NewSelect[string]().
			Key("installType").
			Title("Select Install Type").
			Options(huh.NewOptions("Default", "Custom")...).
			Value(&selectedInstall),

		huh.NewMultiSelect[string]().
			Key("plugins").
			Title("Select Plugins").
			Value(&selectedPlugins).
			OptionsFunc(func() []huh.Option[string] {
				switch agent {
				case "Telegraf":
					switch selectedInstall {
					case "Custom":
						plugins, err := config.LoadPlugins()
						if err != nil {
							return nil
						}
						return huh.NewOptions(plugins.Plugins...)
					default:
						return []huh.Option[string]{
							huh.NewOption(
								"Install: CPU, Disk, Diskio, Kernel, Mem, Processes, Swap, System", "default",
							).Selected(true)}
					}
				default:
					return nil

				}
			}, &selectedInstall),
	)

	return actionGroup
}

func updateAPIKeyGroup(header, apikey, path, defaultPath string) *huh.Group {
	actionGroup := huh.NewGroup(
		huh.NewNote().
			Title(header),

		huh.NewInput().
			Key("apikey").
			Title("Enter your new Hosted Graphite API KEY").
			Prompt("API Key: ").
			Validate(func(s string) error {
				err := utils.ValidateAPIKey(apikey)
				if err != nil {
					return err
				}
				return nil
			}).
			Value(&apikey).
			EchoMode(huh.EchoModePassword),

		huh.NewInput().
			Title("Enter the file path").
			Key("path").
			Prompt("Path: ").
			Description("The default location is already populated. If the path is different please update below.").
			Placeholder(defaultPath).
			Validate(func(s string) error {
				err := telegraf.ValidateFilePath(s)
				if err != nil {
					return err
				}
				return nil
			}).
			Value(&path),
	)

	return actionGroup
}

func uninstallGroup(header string, confirmUninstall bool) *huh.Group {
	actionGroup := huh.NewGroup(
		huh.NewNote().
			Title(header),
		huh.NewConfirm().
			Key("confirm").
			Title("Are you sure you want to uninstall?").
			Description("This will remove the agent and all its configurations.").
			Value(&confirmUninstall),
	)

	return actionGroup
}
