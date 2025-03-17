package pipes

import (
	"os/exec"

	"github.com/hostedgraphite/hg-cli/pipeline"
)

func BrewInstallPipes() []*pipeline.Pipe {
	return []*pipeline.Pipe{
		{
			Name: "Installing Telegraf Agent",
			Cmd:  exec.Command("brew", "install", "telegraf"),
		},
	}
}

func BrewUninstallPipes() []*pipeline.Pipe {
	return []*pipeline.Pipe{
		{
			Name: "Stopping Telegraf Brew Service",
			Cmd:  exec.Command("brew", "services", "stop", "telegraf"),
		},
		{
			Name: "Uninstalling Telegraf Agent",
			Cmd:  exec.Command("brew", "uninstall", "telegraf"),
		},
	}
}
