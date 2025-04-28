package formatters

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/charmbracelet/lipgloss"
	"github.com/hostedgraphite/hg-cli/styles"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var otelSummaryTemplate = `
{{if eq .Action "Install"}}
	{{.SuccessMessage}}
	{{.Config}}
	{{.StartCmd}}
{{else if eq .Action "Update Api Key"}}
	{{.SuccessMessage}}
	{{.Config}}
	{{.RestartCmd}}
{{end}}
`

var telegrafSummaryTemplate = `
{{if eq .Action "Install"}}
	{{.SuccessMessage}}
	{{.Plugins}}
	{{.Config}}
	{{.StartCmd}}
{{else if eq .Action "Update Api Key"}}
	{{.SuccessMessage}}
	{{.Config}}
	{{.RestartCmd}}
{{end}}
`

var (
	titleCaser    = cases.Title(language.English)
	labelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#20b9f7")).Bold(true)
	restartLabel  = labelStyle.Render("Restart Command   : ")
	startLabel    = labelStyle.Render("Start Command     : ")
	configLabel   = labelStyle.Render("Config Path       : ")
	pluginsLabel  = labelStyle.Render("Plugins Installed : ")
	receiverLabel = labelStyle.Render("Receiver          : ")
	exporterLabel = labelStyle.Render("Exporter          : ")
)

var defaultCallToAction = `
To view your metrics, head back to your account, and a dashboard will be automatically added shortly.
For more information on using the hg-cli, check our documentation:
https://docs.hostedgraphite.com/hg-cli
`

var uninstallCallToAction = `
Thanks for using the Hosted Graphite CLI!
The agent has been uninstalled, the hg-cli remains available to assist with your monitoring needs.
`

type ActionSummary struct {
	Agent      string
	Success    bool
	Action     string
	Config     string
	StartCmd   string
	RestartCmd string
	Error      string
}

type OtelContribSummary struct {
	ActionSummary
	Receiver string
	Exporter string
}

type TelegrafSummary struct {
	ActionSummary
	Plugins []string
}

type SummaryContent interface {
	GenerateContent() map[string]string
}

func (o *OtelContribSummary) GenerateContent() map[string]string {
	data := make(map[string]string)

	switch o.ActionSummary.Action {
	case "Install":
		data["StartCmd"] = o.ActionSummary.StartCmd
		data["Config"] = o.ActionSummary.Config
		data["Receiver"] = o.Receiver
		data["Exporter"] = o.Exporter
	case "Update Api Key":
		data["RestartCmd"] = o.ActionSummary.RestartCmd
		data["Config"] = o.ActionSummary.Config
	}
	data["Action"] = o.Action
	data["Agent"] = o.Agent
	data["SuccessMessage"] = o.Action

	return data

}

func (t *TelegrafSummary) GenerateContent() map[string]string {
	data := make(map[string]string)
	switch t.Action {
	case "Install":
		data["StartCmd"] = t.StartCmd
		data["Config"] = t.Config
		data["Plugins"] = strings.Join(t.Plugins, ", ")
	case "Update Api Key":
		data["RestartCmd"] = t.RestartCmd
		data["Config"] = t.Config
	}
	data["Action"] = t.Action
	data["Agent"] = t.Agent
	data["SuccessMessage"] = t.Action

	return data
}

func formatField(key, value string, s styles.Summary) string {
	switch key {
	case "StartCmd":
		return s.Base.Render(s.KeyWord.Render("Start Cmd: ") + s.Items.Render(value))
	case "RestartCmd":
		return s.Base.Render(s.KeyWord.Render("Restart Cmd: ") + s.Items.Render(value))
	case "Config":
		return s.Base.Render(s.KeyWord.Render("Config Location: ") + s.Items.Render(value))
	case "Plugins":
		return s.Base.Render(s.KeyWord.Render("Plugins: ") + s.Items.Render(value))
	case "Receiver":
		return s.Base.Render(s.KeyWord.Render("Receiver: ") + s.Items.Render(value))
	case "Exporter":
		return s.Base.Render(s.KeyWord.Render("Exporter: ") + s.Items.Render(value))
	case "Error":
		return s.Status.Render(fmt.Sprintf("Error: %s", value))
	case "SuccessMessage":
		return s.Status.Render(fmt.Sprintf("Success! - %s complete", value))
	default:
		return value
	}
}

func GenerateSummary(summary SummaryContent, width, height int) string {
	var viewStr strings.Builder
	var err error
	var tmpl *template.Template
	data := summary.GenerateContent()
	s := styles.SummaryStyles(true)
	agent := data["Agent"]

	for key, value := range data {
		data[key] = formatField(key, value, s)
	}

	switch agent {
	case "Telegraf":
		tmpl, err = template.New("telegraf").Parse(telegrafSummaryTemplate)
		if err != nil {
			return fmt.Sprintf("Error parsing template: %v", err)
		}
	case "OpenTelemetry":
		tmpl, err = template.New("otel").Parse(otelSummaryTemplate)
		if err != nil {
			return fmt.Sprintf("Error parsing template: %v", err)
		}
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}
	footer := s.Footer.Render("Thank you for using MetricFire! 🔥 \n(Press q or ctrl+c to quit)")

	content := buf.String()
	viewStr.WriteString(s.Action.Render(styles.MetricfireLogo))
	viewStr.WriteString(content + "\n")
	viewStr.WriteString(s.CtoAction.Render(renderCallToAction(data["Action"], s)) + "\n")
	viewStr.WriteString(footer + "\n")
	content = s.Container.Render(viewStr.String())

	return styles.PlaceContent(width, height, content)
}

func renderCallToAction(action string, s styles.Summary) string {
	var ctoAction string
	switch action {
	case "Install", "Update Api Key":
		ctoAction = defaultCallToAction
	case "Uninstall":
		ctoAction = uninstallCallToAction
	}
	return s.CtoAction.Render(ctoAction)
}

func GenerateCliSummary(summary SummaryContent) string {
	var viewStr strings.Builder
	var cmd, ctoAction, extrasOptions string
	data := summary.GenerateContent()
	agent := data["Agent"]
	action := data["Action"]
	restartCmd := data["RestartCmd"]
	startCmd := data["StartCmd"]
	configPath := data["Config"]

	pipelineTitle := lipgloss.NewStyle().BorderStyle(lipgloss.DoubleBorder()).Width(40).BorderBottom(true).BorderForeground(lipgloss.Color("#f66c00")).Bold(true)

	switch action {
	case "Update Api Key":
		cmd = fmt.Sprintf("%s %s\n", restartLabel, restartCmd)
		ctoAction = defaultCallToAction
	case "Install":
		if agent == "telegraf" {
			extrasOptions = fmt.Sprintf("%s %s\n", pluginsLabel, data["Plugins"])
		} else if agent == "otel" {
			extrasOptions = fmt.Sprintf("%s %s\n%s %s\n", receiverLabel, data["Receiver"], exporterLabel, data["Exporter"])
		}
		cmd = fmt.Sprintf("%s %s\n", startLabel, startCmd)
		ctoAction = defaultCallToAction
	case "Uninstall":
		ctoAction = uninstallCallToAction
		viewStr.WriteString(ctoAction)
		return viewStr.String()
	}

	header := "\n" + titleCaser.String(agent) + " Service Details"

	viewStr.WriteString(pipelineTitle.Render(header))
	viewStr.WriteString("\n" + cmd)
	viewStr.WriteString(fmt.Sprintf("%s %s\n", configLabel, configPath))
	viewStr.WriteString(extrasOptions)
	viewStr.WriteString(ctoAction)

	return viewStr.String()
}
