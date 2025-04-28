package pipes

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/pipeline"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

func DarwinInstallPipes(sysInfo sysinfo.SysInfo) []*pipeline.Pipe {
	var pipes []*pipeline.Pipe
	arch := sysInfo.Arch

	latest, err := utils.GetLatestReleaseTag("open-telemetry", "opentelemetry-collector-releases")
	if err != nil {
		latest = "v0.123.1" // Default
	}
	release := fmt.Sprintf("otelcol-contrib_%s_darwin_%s.tar.gz", latest[1:], arch)
	url := "https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/" + latest + "/" + release
	tmpDir := "/tmp/hg-cli"

	pipes = []*pipeline.Pipe{
		{
			Name: "Creating Temporary Dir",
			Cmd: exec.Command(
				"mkdir",
				"-p",
				tmpDir,
			),
		},
		{
			Name: "Downloading OpenTelemetry to " + tmpDir,
			Cmd: exec.Command(
				"curl",
				"--tlsv1.2",
				"-fL",
				"-o",
				tmpDir+release,
				url,
			),
		},
		{
			Name: "Starting Extraction of Tar Files",
			Cmd: exec.Command(
				"tar",
				"-xvf",
				tmpDir+release,
				"-C",
				"/tmp/hg-cli",
			),
		},
		{
			Name: "Moving Exe File to /usr/local/bin",
			Cmd: exec.Command(
				"sh",
				"-c",
				`mkdir -p /usr/local/bin && mv /tmp/hg-cli/otelcol-contrib /usr/local/bin/`,
			),
		},
		{
			Name: "Cleaning up Temporary Directory",
			Cmd: exec.Command(
				"rm",
				"-rf",
				tmpDir,
			),
		},
	}

	return pipes
}

func DarwinConfigPipes(options map[string]interface{}, serviceSettings map[string]string, plistFile string) []*pipeline.Pipe {
	plistPath := "/usr/local/etc/otelcol-contrib/com.otelcol-contrib-agent.plist"
	homeDir := os.Getenv("HOME")
	plistDest := homeDir + "/Library/LaunchAgents/com.otelcol-contrib-agent.plist"

	pipes := []*pipeline.Pipe{
		{
			Name: "Creating Config.Yaml",
			Cmd: exec.Command(
				"sh",
				"-c",
				`mkdir -p /usr/local/etc/otelcol-contrib && touch /usr/local/etc/otelcol-contrib/config.yaml`,
			),
		},
		{
			Name: "Creating Plist File",
			Cmd: exec.Command(
				"touch",
				"/usr/local/etc/otelcol-contrib/com.otelcol-contrib-agent.plist",
			),
		},
		{
			Name: "Updating com.otelcocom.otelcol-contrib-agent.plistl-contrib-agent.plist",
			Cmd: exec.Command(
				"bash",
				"-c",
				fmt.Sprintf("echo '%s' > /usr/local/etc/otelcol-contrib/com.otelcol-contrib-agent.plist", plistFile),
			),
		},
		{
			Name: "Moving Plist File to Launch Daemons",
			Cmd: exec.Command(
				"mv",
				plistPath,
				plistDest,
			),
		},
	}

	return pipes
}

func DarwinUninstallPipes() []*pipeline.Pipe {
	// otel is installed and used on a user level on mac
	// since the cli command needs sudo to be ran
	// it causes issues during stop/unloading.
	// The orignial user not sudo is used and passed to the
	// launch commands.
	origUser := os.Getenv("SUDO_USER")
	plistPath := fmt.Sprintf("/Users/%s/Library/LaunchAgents/com.otelcol-contrib-agent.plist", origUser)

	pipes := []*pipeline.Pipe{
		{
			Name: "Stopping Otel-Contrib Agent",
			Cmd: exec.Command(
				"sudo",
				"-u",
				origUser,
				"launchctl",
				"stop",
				"com.otelcol-contrib-agent",
			),
		},
		{
			Name: "Unloading Otel-Contrib Agent",
			Cmd: exec.Command(
				"sudo",
				"-u",
				origUser,
				"launchctl",
				"unload",
				plistPath,
			),
		},
		{
			Name: "Removing Otel-Contrib Agent Plist File",
			Cmd: exec.Command(
				"sh",
				"-c",
				`rm ~/Library/LaunchAgents/com.otelcol-contrib-agent.plist`,
			),
		},
		{
			Name: "Removing Otel-Contrib Binary",
			Cmd: exec.Command(
				"rm",
				"/usr/local/bin/otelcol-contrib",
			),
		},
	}

	return pipes
}
