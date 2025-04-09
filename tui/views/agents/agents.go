package agents

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf"
	"github.com/hostedgraphite/hg-cli/styles"
	"github.com/hostedgraphite/hg-cli/tui/views/config"
	"github.com/hostedgraphite/hg-cli/utils"
)

type Telegraf struct {
	apikey           string
	selectedInstall  string
	selectedPlugins  []string
	confirmUninstall bool
	path             string
	header           string
}

func (t *Telegraf) InstallView(apikey, selectedInstall string, selectedPlugins []string) (*huh.Group, error) {
	installGroup := huh.NewGroup(
		huh.NewNote().
			Title(t.header),

		huh.NewInput().
			Key("apikey").
			Title("Enter your Hosted Graphite API key").
			Prompt("API Key: ").
			Validate(func(s string) error {
				err := utils.ValidateAPIKey(t.apikey)
				if err != nil {
					return err
				}
				return nil
			}).
			Value(&t.apikey).
			EchoMode(huh.EchoModePassword),

		huh.NewSelect[string]().
			Key("installType").
			Title("Select Install Type").
			Options(huh.NewOptions("Default", "Custom")...).
			Value(&t.selectedInstall),

		huh.NewMultiSelect[string]().
			Key("plugins").
			Title("Select Plugins").
			Value(&t.selectedPlugins).
			OptionsFunc(func() []huh.Option[string] {
				switch t.selectedInstall {
				case "Custom":
					plugins, err := config.LoadPlugins()
					if err != nil {
						return nil
					}
					return huh.NewOptions(plugins.Plugins...)
				default:
					defaultPlugins := telegraf.DefaultTelegrafPlugins
					options := huh.NewOptions(defaultPlugins...)
					for i := range options {
						options[i] = options[i].Selected(true)
					}
					return options
				}
			}, &t.selectedInstall),
	)

	return installGroup, nil
}
func (t *Telegraf) UninstallView(confirm bool) (*huh.Group, error) {
	uninstallGroup := huh.NewGroup(
		huh.NewNote().
			Title(t.header),
		huh.NewConfirm().
			Key("confirmUninstall").
			Title("Are you sure you want to uninstall Telegraf?").
			Description("This will remove the agent, but not the configuration files").
			Value(&t.confirmUninstall),
	)

	return uninstallGroup, nil
}
func (t *Telegraf) UpdateApiKeyView(apikey, path, defaultPath string) (*huh.Group, error) {

	updateGroup := huh.NewGroup(
		huh.NewNote().
			Title(t.header),

		huh.NewInput().
			Key("apikey").
			Title("Enter your new Hosted Graphite API key").
			Prompt("API Key: ").
			Validate(func(s string) error {
				err := utils.ValidateAPIKey(t.apikey)
				if err != nil {
					return err
				}
				return nil
			}).
			Value(&t.apikey).
			EchoMode(huh.EchoModePassword),

		huh.NewInput().
			Key("path").
			Title("Enter the path to the Telegraf configuration file").
			Prompt("Path: ").
			Description("The default location is already populated. If the path is different please update below.").
			Placeholder(defaultPath).
			Value(&t.path).
			Validate(func(s string) error {
				if s == "" {
					s = defaultPath
				}
				err := telegraf.ValidateFilePath(s)
				if err != nil {
					return err
				}
				return nil
			}),
	)

	return updateGroup, nil
}

type Otel struct{}

func (o *Otel) InstallView()      {}
func (o *Otel) UninstallView()    {}
func (o *Otel) UpdateApiKeyView() {}

func NewAgentsFields(agent string) AgentsFieldViews {
	header := getHeader(agent)
	switch strings.ToLower(agent) {
	case "telegraf":
		return &Telegraf{
			header: header,
		}
	case "otel":
		return nil
	default:
		return nil
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
