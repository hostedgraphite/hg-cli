package installers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
)

var installSteps = []string{
	"Step 1. Downloading Telegraf to ~\\Downloads",
	"Step 2. Expanding Archive to C:\\Program Files",
	"Step 3. Moving exe to C:\\Program File\\InfluxData\\telegraf",
	"Step 4. Installing Telegraf as windows service",
}

var installStepsTmpls = []string{
	`$ProgressPreference='SilentlyContinue';Invoke-WebRequest -Uri https://dl.influxdata.com/telegraf/releases/{{.Release}} -OutFile ~\Downloads\{{.Release}};`,
	`$ProgressPreference='SilentlyContinue';Expand-Archive ~\Downloads\{{.Release}} -DestinationPath 'C:\Program Files\InfluxData\telegraf\'`,
	`Move-Item "C:\Program Files\InfluxData\telegraf\telegraf-{{.Latest}}\telegraf.*" "C:\Program Files\InfluxData\telegraf\"`,
	`& "C:\Program Files\InfluxData\telegraf\telegraf.exe" --config "C:\Program Files\InfluxData\telegraf\telegraf.conf" --service-name telegraf service install`,
}

type installContext struct {
	Latest  string
	Release string
}

func TelegrafAgentInstallWindows(arch string, updates chan<- string) error {

	if IsInstalledWindows() {
		return fmt.Errorf("telegraf is already installed at C:\\Program Files\\InfluxData\\telegraf\nAborting installation")
	}

	var err error
	latest, err := utils.GetLatestReleaseTag("influxdata", "telegraf")
	if err != nil {
		return fmt.Errorf("error getting latest telegraf release: %v", err)
	}

	latest = latest[1:]
	release := fmt.Sprintf("telegraf-%s_windows_%s.zip", latest, arch)

	context := installContext{
		Latest:  latest,
		Release: release,
	}

	commandSteps, err := buildCommandSteps(context)
	if err != nil {
		return err
	}

	for _, stepname := range installSteps {
		updates <- stepname
		command := commandSteps[stepname]
		shell := determineShell()

		cmd := exec.Command(shell, "-Command", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error installing telegraf: %v\nOutput: %s", err, output)
		}
	}

	return err
}

func determineShell() string {
	shell := "pwsh" // powershell core 7
	// check if pwsh.exe exits
	if _, err := exec.LookPath(shell); err != nil {
		shell = "powershell" // PowerShell 5
	}
	return shell
}

func buildCommandSteps(context installContext) (map[string]string, error) {

	commandSteps := make(map[string]string)

	for index, commandtmpl := range installStepsTmpls {
		stepname := installSteps[index]

		tmpl, err := template.New("installStep").Parse(commandtmpl)
		if err != nil {
			return nil, fmt.Errorf("unable to build tempalte for install step: %s \nError: %v", stepname, err)
		}

		var rendered_cmd bytes.Buffer
		err = tmpl.Execute(&rendered_cmd, context)
		if err != nil {
			return nil, fmt.Errorf("unable to render install step: %s \nError: %v", stepname, err)
		}

		commandSteps[stepname] = rendered_cmd.String()
	}

	return commandSteps, nil
}

func checkFileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsInstalledWindows() bool {
	// Check if specific files exist
	filesToCheck := []string{
		"C:\\Program Files\\InfluxData\\telegraf\\telegraf.exe",
		"C:\\Program Files\\InfluxData\\telegraf\\telegraf.conf",
	}

	for _, file := range filesToCheck {
		if checkFileExists(file) {
			return true
		}
	}

	return false
}
