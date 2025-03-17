package pipes

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/pipeline"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

func WindowsInstallPipes(sysInfo sysinfo.SysInfo) ([]*pipeline.Pipe, error) {

	if IsInstalledWindows() {
		return nil, fmt.Errorf("telegraf is already installed. Please check C:\\Program Files\\InfluxData\\telegraf")
	}

	latest, err := utils.GetLatestReleaseTag("influxdata", "telegraf")
	if err != nil {
		latest = "v1.33.1" // Default
	}

	latest = latest[1:]
	arch := sysInfo.Arch
	release := fmt.Sprintf("telegraf-%s_windows_%s.zip", latest, arch)
	shell := determineShell()

	pipes := []*pipeline.Pipe{
		{
			Name: "Downloading telegraf to ~\\Downloads",
			Cmd:  exec.Command(shell, "-Command", `$ProgressPreference='SilentlyContinue';Invoke-WebRequest -Uri https://dl.influxdata.com/telegraf/releases/`+release+` -OutFile ~\Downloads\`+release+`;`),
		},
		{
			Name: "Expanding telegraf archive to C:\\Program Files",
			Cmd:  exec.Command(shell, "-Command", `$ProgressPreference='SilentlyContinue';Expand-Archive ~\Downloads\`+release+` -DestinationPath 'C:\Program Files\InfluxData\telegraf\'`),
		},
		{
			Name: "Moving telegraf exe to C:\\Program File\\InfluxData\\telegraf",
			Cmd:  exec.Command(shell, "-Command", `Move-Item "C:\Program Files\InfluxData\telegraf\telegraf-`+latest+`\telegraf.*" "C:\Program Files\InfluxData\telegraf\"`),
		},
		{
			Name: "Installing Telegraf as Windows service",
			Cmd:  exec.Command(shell, "-Command", `& "C:\Program Files\InfluxData\telegraf\telegraf.exe" --service-name telegraf service install`),
		},
	}

	return pipes, nil
}

func WindowsConfigPipes(options map[string]interface{}, serviceSettings map[string]string) []*pipeline.Pipe {
	inputs := strings.Join(options["plugins"].([]string), ":")
	telegrafCmd := serviceSettings["serviceCmd"]
	configpath := serviceSettings["configPath"]

	pipes := []*pipeline.Pipe{
		pipeline.NewPipe("Configuring Telegraf Plugins", exec.Command("powershell", "-Command", fmt.Sprintf("& '%s' --input-filter %s --output-filter graphite config", telegrafCmd, inputs))).PostRun(
			func(ctx context.Context) error {
				output := ctx.Value("output").(string)
				err := os.WriteFile(configpath, []byte(output), 0644)
				return err
			},
		),
	}

	return pipes
}

func determineShell() string {
	shell := "pwsh" // powershell core 7
	// check if pwsh.exe exits
	if _, err := exec.LookPath(shell); err != nil {
		shell = "powershell" // PowerShell 5
	}
	return shell
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
