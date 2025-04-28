package pipes

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/pipeline"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

func WindowsInstallPipes(sysInfo sysinfo.SysInfo) ([]*pipeline.Pipe, error) {
	if isInstalledWindows() {
		return nil, fmt.Errorf("otelcontribcol is already installed. Please check C:\\Program Files\\OpenTelemetry Collector Contrib")
	}

	latest, err := utils.GetLatestReleaseTag("open-telemetry", "opentelemetry-collector-releases")
	if err != nil {
		latest = "v0.123.1" // Default
	}

	arch := sysInfo.Arch
	release := fmt.Sprintf("otelcol-contrib_%s_windows_%s.tar.gz", latest[1:], arch)
	shell := determineShell()
	uri := "https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/" + latest + "/" + release

	pipes := []*pipeline.Pipe{
		{
			Name: "Downloading otelcontribcol to ~\\Downloads",
			Cmd:  exec.Command(shell, "-Command", `$ProgressPreference='SilentlyContinue';Invoke-WebRequest -Uri `+uri+` -OutFile $env:USERPROFILE\Downloads\`+release+`;`),
		},
		{
			Name: "Creating directory for otelcontribcol extraction",
			Cmd:  exec.Command(shell, "-Command", `New-Item -ItemType Directory -Path 'C:\Program Files\OpenTelemetry Collector Contrib'`),
		},
		{
			Name: "Expanding otelcontribcol archive to C:\\Program Files\\OpenTelemetry Collector Contrib",
			Cmd:  exec.Command(shell, "-Command", `$ProgressPreference='SilentlyContinue';tar -xzf $env:USERPROFILE\Downloads\`+release+` -C "C:\\Program Files\\OpenTelemetry Collector Contrib"`),
		},
	}
	return pipes, nil

}

func WindowsConfigPipes(options map[string]interface{}, serviceSettings map[string]string) []*pipeline.Pipe {
	shell := determineShell()
	configPath := serviceSettings["configPath"]
	exePath := serviceSettings["exePath"]

	binPathName := fmt.Sprintf(`'"%s" --config "%s"'`, exePath, configPath)

	pipes := []*pipeline.Pipe{
		{
			Name: "Creating Configuration file",
			Cmd:  exec.Command(shell, "-Command", fmt.Sprintf(`New-Item -ItemType File -Path '%s'`, configPath)),
		},
		// kind of wierd place to put, but the config file is needed for the service to be created
		{
			Name: "Creating OpenTelemetry Service",
			Cmd: exec.Command(
				shell,
				"-Command",
				"New-Service",
				"-Name",
				"otelcol-contrib",
				"-BinaryPathName",
				binPathName,
			),
		},
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

func isInstalledWindows() bool {
	filesToCheck := []string{
		"C:\\Program Files\\OpenTelemetry Collector Contrib\\otelcontribcol.exe",
		"C:\\Program Files\\OpenTelemetry Collector Contrib\\config.yaml",
	}

	for _, file := range filesToCheck {
		if checkFileExists(file) {
			return true
		}
	}

	return false
}

func WindowsUninstallPipes(sysInfo sysinfo.SysInfo) ([]*pipeline.Pipe, error) {
	shell := determineShell()
	if !isInstalledWindows() {
		return nil, fmt.Errorf("no exe found - Unable to remove service")
	}

	pipes := []*pipeline.Pipe{
		{
			Name: "Stopping otelcontribcol service",
			Cmd: exec.Command(
				shell,
				"-Command",
				"sc.exe",
				"stop",
				"otelcol-contrib",
			),
		},
		{
			Name: "Uninstalling otelcontribcol service",
			Cmd: exec.Command(
				shell,
				"-Command",
				"sc.exe",
				"delete",
				"otelcol-contrib",
			),
		},
	}
	return pipes, nil
}
