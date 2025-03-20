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

var (
	titleCaser   = cases.Title(language.English)
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#20b9f7")).Bold(true)
	restartLabel = labelStyle.Render("Restart Command   : ")
	startLabel   = labelStyle.Render("Start Command     : ")
	configLabel  = labelStyle.Render("Config Path       : ")
	pluginsLabel = labelStyle.Render("Plugins Installed : ")
)

var defaultCallToAction = `
To view your metrics, head back to your account, and a dashboard will be automatically added shortly.
For more information on using the hg-cli, visit the documentation:
https://docs.hostedgraphite.com/hg-cli
`

var uninstallCallToAction = `
Thanks for running the Hosted Graphite CLI—this has now been uninstalled. Here to help with your next monitoring step!
`

type ActionSummary struct {
	Agent      string
	Success    bool
	Action     string
	Plugins    []string
	Config     string
	StartCmd   string
	RestartCmd string
	Error      string
}

func GenerateSummary(action ActionSummary, width, height int) string {
	var viewStr strings.Builder
	var summary, title, ctoaction string

	s := styles.SummaryStyles(action.Success)
	switch action.Action {
	case "Install":
		title = "Install Agent"
		ctoaction = defaultCallToAction
	case "Update Api Key":
		title = "Update Api Key"
		ctoaction = defaultCallToAction
	case "Uninstall":
		title = "Uninstall Agent"
		ctoaction = uninstallCallToAction
	}

	footer := s.Footer.Render("Thank you for using MetriFire! 🔥 \n(Press q or ctrl+c to quit)")

	tmpl, err := template.New("summary").Parse(`
{{.ActionTitle}}
{{if eq .Action "Install"}}
    {{if .Success}}
        {{.SuccessMessage}}
        {{.Plugins}}
        {{.Config}}
        {{.StartCmd}}
    {{else}}
        {{.FailureMessage}}
        {{.Error}}
    {{end}}
{{else if eq .Action "Update Api Key"}}
    {{if .Success}}
        {{.SuccessMessage}}
        {{.Config}}
        {{.RestartCmd}}
    {{else}}
        {{.FailureMessage}}
        {{.Error}}
    {{end}}
{{end}}
`)
	if err != nil {
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	data := map[string]interface{}{
		"ActionTitle":    s.Action.Render(title),
		"Action":         action.Action,
		"Success":        action.Success,
		"SuccessMessage": s.Status.Render(fmt.Sprintf("Success! - %s did %s successfully", action.Agent, action.Action)),
		"Plugins":        s.Base.Render(s.KeyWord.Render("Plugins: ") + s.Items.Render(strings.Join(action.Plugins, ", "))),
		"Config":         s.Base.Render(s.KeyWord.Render("Config Location: ") + s.Items.Render(action.Config)),
		"StartCmd":       s.Base.Render(s.KeyWord.Render("Start Cmd: ") + s.Items.Render(action.StartCmd)),
		"RestartCmd":     s.Base.Render(s.KeyWord.Render("Restart Cmd: ") + s.Items.Render(action.RestartCmd)),
		"FailureMessage": s.Status.Render(fmt.Sprintf("Failure - %s did not %s successfully\n", action.Agent, action.Action)),
		"Error":          s.Base.Render(s.KeyWord.Render("Errors: ") + s.Items.Render(action.Error)),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}

	summary = buf.String()

	viewStr.WriteString(s.Action.Render(styles.MetricfireLogo))
	viewStr.WriteString(summary + "\n")
	viewStr.WriteString(s.CtoAction.Render(ctoaction) + "\n")
	viewStr.WriteString(footer + "\n")
	content := s.Container.Render(viewStr.String())

	return styles.PlaceContent(width, height, content)
}

func GenerateCliSummary(action ActionSummary) string {
	var viewStr strings.Builder
	var plugins, cmd, ctoAction string

	pipelineTitle := lipgloss.NewStyle().BorderStyle(lipgloss.DoubleBorder()).Width(40).BorderBottom(true).BorderForeground(lipgloss.Color("#f66c00")).Bold(true)

	switch action.Action {
	case "Update Api Key":
		cmd = fmt.Sprintf("%s %s\n", restartLabel, action.RestartCmd)
		ctoAction = defaultCallToAction
	case "Install":
		plugins = fmt.Sprintf("%s %s\n", pluginsLabel, strings.Join(action.Plugins, ", "))
		cmd = fmt.Sprintf("%s %s\n", startLabel, action.StartCmd)
		ctoAction = defaultCallToAction
	case "Uninstall":
		ctoAction = uninstallCallToAction
	}

	header := "\n" + titleCaser.String(action.Agent) + " Service Details"

	viewStr.WriteString(pipelineTitle.Render(header))
	viewStr.WriteString("\n" + cmd)
	viewStr.WriteString(fmt.Sprintf("%s %s\n", configLabel, action.Config))
	viewStr.WriteString(plugins)
	viewStr.WriteString(ctoAction)

	return viewStr.String()
}
