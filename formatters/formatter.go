package formatters

import (
	"bytes"
	"fmt"
	"hg-cli/styles"
	"strings"
	"text/template"
)

type ActionSummary struct {
	Agent    string
	Success  bool
	Action   string
	Plugins  []string
	Config   string
	StartCmd string
	Error    string
}

func GenerateSummary(action ActionSummary, width, height int) string {
	var viewStr strings.Builder
	var summary, title string

	s := styles.SummaryStyles(action.Success)
	switch action.Action {
	case "Install":
		title = "Install Agent"
	case "Update Api Key":
		title = "Update Api Key"
	case "Uninstall":
		title = "Uninstall Agent"
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
        {{.StartCmd}}
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
	viewStr.WriteString(footer + "\n")
	content := s.Container.Render(viewStr.String())

	return styles.PlaceContent(width, height, content)
}

func GenerateCliSummary(action ActionSummary) string {
	var viewStr strings.Builder
	var summary, title, plugins, cmd string
	s := styles.DefaultStyles()

	switch action.Action {
	case "Update Api Key":
		title = "Update Api Key Summary"
		plugins = ""
		cmd = fmt.Sprintf("Restart cmd: %s\n", action.StartCmd)
	case "Install":
		title = "Installation Summary"
		plugins = fmt.Sprintf("Plugins: %s\n", strings.Join(action.Plugins, ", "))
		cmd = fmt.Sprintf("Start cmd: %s\n", action.StartCmd)
	}

	viewStr.WriteString("\n" + title + "\n")
	viewStr.WriteString(fmt.Sprintf("Agent: %s\n", action.Agent))
	viewStr.WriteString(cmd)
	viewStr.WriteString(fmt.Sprintf("Config: %s\n", action.Config))
	viewStr.WriteString(plugins)

	summary = s.Cli.Render(viewStr.String())

	return summary
}
