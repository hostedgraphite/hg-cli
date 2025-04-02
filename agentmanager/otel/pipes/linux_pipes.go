package pipes

import (
	"os/exec"

	"github.com/hostedgraphite/hg-cli/pipeline"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

func LinuxInstallPipes(sysInfo sysinfo.SysInfo) []*pipeline.Pipe {
	var pipes []*pipeline.Pipe
	pkgMngr := sysInfo.PkgMngr

	if pkgMngr == "brew" {
		return pipes
	} else if pkgMngr == "apt" {
		pipes = aptInstallPipes()
	} else if pkgMngr == "yum" || pkgMngr == "dnf" {
		pipes = yumInstallPipes()
		return pipes
	} else {
		return pipes
	}

	return pipes
}

func aptInstallPipes() []*pipeline.Pipe {
	packagePath := "https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v0.118.0/otelcol-contrib_0.118.0_linux_amd64.deb"

	tmpDir := "/tmp/hg-cli"
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
			Name: "Installing Otel-Contrib",
			Cmd:  exec.Command("dpkg", "-i", tmpDir+"/otelcol-contrib_0.118.0_linux_amd64.deb"),
		},
	}
	return pipes
}

func yumInstallPipes() []*pipeline.Pipe {
	packagePath := "https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v0.119.0/otelcol-contrib_0.119.0_linux_amd64.rpm"

	tmpDir := "/tmp/hg-cli"
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
			Name: "Installing Otel-Contrib",
			Cmd:  exec.Command("rpm", "-ivh", tmpDir+"/otelcol-contrib_0.119.0_linux_amd64.rpm"),
		},
	}
	return pipes
}
