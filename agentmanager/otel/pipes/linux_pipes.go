package pipes

import (
	"fmt"
	"os/exec"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/pipeline"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

func LinuxInstallPipes(sysInfo sysinfo.SysInfo) []*pipeline.Pipe {
	var pipes []*pipeline.Pipe
	pkgMngr := sysInfo.PkgMngr
	arch := sysInfo.Arch

	latest, err := utils.GetLatestReleaseTag("open-telemetry", "opentelemetry-collector-releases")
	if err != nil {
		latest = "v0.123.1" // Default
	}
	release := fmt.Sprintf("otelcol-contrib_%s_linux_%s", latest[1:], arch)
	packagePath := "https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/" + latest + "/" + release

	if pkgMngr == "apt" {
		pipes = aptInstallPipes(packagePath, release)
	} else if pkgMngr == "yum" || pkgMngr == "dnf" {
		pipes = yumInstallPipes(packagePath, release)
	} else {
		pipes = manualInstallPipes(packagePath, release)
	}

	return pipes
}

func LinuxManualConfigPipes(options map[string]interface{}, serviceSettings map[string]string, sytemdFile string) []*pipeline.Pipe {
	pipes := []*pipeline.Pipe{
		{
			Name: "Creating Otel-Contrib Config Directory",
			Cmd: exec.Command(
				"mkdir",
				"/etc/otelcol-contrib/",
			),
		},
		{
			Name: "Creating Otel-Contrib Config File",
			Cmd: exec.Command(
				"touch",
				"/etc/otelcol-contrib/config.yaml",
			),
		},
		{
			Name: "Creating Otel-Contrib Systemd File",
			Cmd: exec.Command(
				"touch",
				"/etc/systemd/system/otelcol-contrib.service",
			),
		},
		{
			Name: "Creating Otel-Contrib Systemd File",
			Cmd: exec.Command(
				"bash",
				"-c",
				fmt.Sprintf("echo '%s' > /etc/systemd/system/otelcol-contrib.service", sytemdFile),
			),
		},
	}
	return pipes
}

func linuxManualUninstallPipes() []*pipeline.Pipe {
	pipes := []*pipeline.Pipe{
		{
			Name: "Stopping Otel-Contrib Service",
			Cmd: exec.Command(
				"systemctl",
				"stop",
				"otelcol-contrib",
			),
		},
		{
			Name: "Uninstalling Otel-Contrib",
			Cmd: exec.Command(
				"rm",
				"-rf",
				"/usr/bin/otelcol-contrib",
			),
		},
		{
			Name: "Removing Otel-Contrib Service",
			Cmd: exec.Command(
				"rm",
				"-rf",
				"/etc/systemd/system/otelcol-contrib.service",
			),
		},
	}
	return pipes
}

func aptInstallPipes(packagePath, release string) []*pipeline.Pipe {
	tmpDir := "/tmp/hg-cli/"
	packagePath = packagePath + ".deb"
	debPath := tmpDir + release + ".deb"
	pipes := []*pipeline.Pipe{
		{
			Name: "Creating TMP Directory",
			Cmd:  exec.Command("mkdir", "-p", tmpDir),
		},
		{
			Name: "Downloading Otel-Contrib Package",
			Cmd:  exec.Command("wget", "-P", tmpDir, packagePath),
		},
		{
			Name: "Installing Otel-Contrib ",
			Cmd:  exec.Command("dpkg", "-i", debPath),
		},
	}
	return pipes
}

func manualInstallPipes(packagePath, release string) []*pipeline.Pipe {
	tmpDir := "/tmp/hg-cli/"
	packagePath = packagePath + ".tar.gz"
	tarPath := tmpDir + release + ".tar.gz"
	pipes := []*pipeline.Pipe{
		{
			Name: "Creating TMP Directory",
			Cmd:  exec.Command("mkdir", "-p", tmpDir),
		},
		{
			Name: "Downloading OpenTelemetry to " + tmpDir,
			Cmd: exec.Command(
				"curl",
				"--tlsv1.2",
				"-fL",
				"-o",
				tarPath,
				packagePath,
			),
		},
		{
			Name: "Starting Extraction of Tar Files",
			Cmd: exec.Command(
				"tar",
				"-xvf",
				tarPath,
				"-C",
				"/tmp/hg-cli",
			),
		},
		{
			Name: "Moving Exe File to /usr/local/bin",
			Cmd: exec.Command(
				"mv",
				"/tmp/hg-cli/otelcol-contrib",
				"/usr/bin/",
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

func yumInstallPipes(packagePath, release string) []*pipeline.Pipe {
	tmpDir := "/tmp/hg-cli/"
	packagePath = packagePath + ".rpm"
	debPath := tmpDir + release + ".rpm"
	pipes := []*pipeline.Pipe{
		{
			Name: "Creating TMP Directory",
			Cmd:  exec.Command("mkdir", "-p", tmpDir),
		},
		{
			Name: "Downloading Otel-Contrib Package",
			Cmd:  exec.Command("wget", "-P", tmpDir, packagePath),
		},
		{
			Name: "Installing Otel-Contrib ",
			Cmd:  exec.Command("rpm", "-ivh", debPath),
		},
	}
	return pipes
}

func LinxUninstallPipes(sysInfo sysinfo.SysInfo) []*pipeline.Pipe {
	var pipes []*pipeline.Pipe
	pkgMngr := sysInfo.PkgMngr

	if pkgMngr == "apt" {
		pipes = linuxDebUninstall()
	} else if pkgMngr == "yum" || pkgMngr == "dnf" {
		pipes = linuxRpmUninstall()
	} else {
		pipes = linuxManualUninstallPipes()
	}

	return pipes
}

func linuxDebUninstall() []*pipeline.Pipe {
	pipes := []*pipeline.Pipe{
		{
			Name: "Uninstalling Otel-Contrib",
			Cmd: exec.Command(
				"dpkg",
				"-r",
				"otelcol-contrib",
			),
		},
	}

	return pipes
}

func linuxRpmUninstall() []*pipeline.Pipe {
	pipes := []*pipeline.Pipe{
		{
			Name: "Uninstalling Otel-Contrib",
			Cmd: exec.Command(
				"rpm",
				"-e",
				"otelcol-contrib",
			),
		},
	}

	return pipes
}
