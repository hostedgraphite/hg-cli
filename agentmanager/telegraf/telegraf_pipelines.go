package telegraf

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	telegrafPipes "github.com/hostedgraphite/hg-cli/agentmanager/telegraf/pipes"
	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/pipeline"
)

func (t *Telegraf) InstallPipeline(updates chan *pipeline.Pipe) (*pipeline.Pipeline, error) {
	var err error
	var sysInfo = t.sysinfo
	var pipes []*pipeline.Pipe

	switch sysInfo.Os {
	case "linux":
		pipes = telegrafPipes.LinuxInstallPipes(sysInfo)
	case "darwin":
		pipes = telegrafPipes.DarwinInstallPipes(sysInfo)
	case "windows":
		pipes, err = telegrafPipes.WindowsInstallPipes(sysInfo)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported operating system: %v", err)
	}

	configPipes, err := t.configPipeline()
	if err != nil {
		return nil, err
	}
	pipes = append(pipes, configPipes...)

	pipeline := pipeline.NewPipeline(fmt.Sprintf("Installing Telegraf Agent (%s-%s)", sysInfo.Os, sysInfo.PkgMngr), pipes, updates)

	return &pipeline, err
}

func (t *Telegraf) configPipeline() ([]*pipeline.Pipe, error) {
	var err error
	// var apikey = t.apikey
	var sysInfo = t.sysinfo
	var options = t.options
	var serviceSettings = t.serviceSettings
	var pipes []*pipeline.Pipe

	switch sysInfo.Os {
	case "windows":
		pipes = telegrafPipes.WindowsConfigPipes(options, serviceSettings)
	default:
		// "Should" work for call linux systems...
		pipes = telegrafPipes.LinuxConfigPipes(options, serviceSettings)
	}

	updatePipe := t.graphiteOutputUpdatePipe()

	pipes = append(pipes, updatePipe...)

	return pipes, err
}

func (t *Telegraf) graphiteOutputUpdatePipe() []*pipeline.Pipe {
	pipes := []*pipeline.Pipe{
		pipeline.NewPipe("Updating Telegraf Graphite Output Config", exec.Command("sleep", "0.2s")).PostRun(
			func(ctx context.Context) error {
				return graphiteOutputUpdate(t.apikey, t.serviceSettings["configPath"])
			},
		),
	}
	return pipes
}

func graphiteOutputUpdate(apikey, configPath string) error {
	fullConfig, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	graphiteBlock := `\[\[outputs\.graphite\]\](?:.|\s)*?\[\[`

	updates := map[string]string{
		`prefix\s*=\s*".*?"`:    fmt.Sprintf(`prefix = "%s.telegraf"`, apikey),
		`servers\s*=\s*\[.*?\]`: `servers = ["carbon.hostedgraphite.com:2003"]`,
		`template\s*=\s*".*?"`:  `## template = "host.tags.measurement.field"`,
	}
	updatedConfig, err := utils.UpdateConfigBlock(string(fullConfig), graphiteBlock, updates)

	if err != nil {
		return fmt.Errorf("error during updating: %v", err)
	}

	err = os.WriteFile(configPath, []byte(updatedConfig), 0644)

	if err != nil {
		return fmt.Errorf("error writing file:%v", err)
	}

	return nil
}
