package telegraf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"

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
	os := t.sysinfo.Os
	var cmd *exec.Cmd

	if os == "windows" {
		cmd = exec.Command("powershell", "-Command", "echo test")
	} else {
		cmd = exec.Command("sleep", "1")
	}

	pipes := []*pipeline.Pipe{
		pipeline.NewPipe("Updating Telegraf Graphite Output Config", cmd).PostRun(
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

func (t *Telegraf) UninstallPipeline(updates chan *pipeline.Pipe) (*pipeline.Pipeline, error) {
	var err error
	var sysInfo = t.sysinfo
	var pipes []*pipeline.Pipe

	switch sysInfo.Os {
	case "linux":
		pipes = telegrafPipes.LinuxUninstallPipes(sysInfo)
	case "darwin":
		pipes = telegrafPipes.DarwinUninstallPipes(sysInfo)
	case "windows":
		pipes, err = telegrafPipes.WindowsUninstallPipes()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported operating system: %v", err)
	}

	pipeline := pipeline.NewPipeline(fmt.Sprintf("Uninstalling Telegraf Agent (%s-%s)", sysInfo.Os, sysInfo.PkgMngr), pipes, updates)

	return &pipeline, err
}

func (t *Telegraf) UpdateApiKeyPipeline(updates chan *pipeline.Pipe) (*pipeline.Pipeline, error) {
	var err error
	var sysInfo = t.sysinfo
	var pipes []*pipeline.Pipe
	var filePath = t.options["configPath"].(string)
	var apikey = t.apikey

	switch sysInfo.Os {
	case "linux":
		pipes = telegrafPipes.LinuxUpdateApiKeyPipe(apikey, filePath)
	case "darwin", "windows":
		pipes = defaultApiKeyPipe(apikey, filePath)
	default:
		return nil, fmt.Errorf("unsupported operating system: %v", err)
	}

	pipeline := pipeline.NewPipeline(fmt.Sprintf("Uninstalling Telegraf Agent (%s-%s)", sysInfo.Os, sysInfo.PkgMngr), pipes, updates)

	return &pipeline, err
}

func defaultApiKeyPipe(apikey, filePath string) []*pipeline.Pipe {
	cmd := exec.Command("sleep", "1")

	pipes := []*pipeline.Pipe{
		pipeline.NewPipe("Updating Telegraf Config", cmd).PostRun(
			func(ctx context.Context) error {
				return apiUpdater(apikey, filePath)
			},
		),
	}

	return pipes
}

func apiUpdater(apikey, filePath string) error {
	newPrefix := apikey + ".telegraf"
	re := regexp.MustCompile(`(?m)^\s*prefix\s*=\s*".*"`)
	currentConfig, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	updatedConfig := re.ReplaceAllString(string(currentConfig), fmt.Sprintf(`prefix = "%s"`, newPrefix))

	err = os.WriteFile(filePath, []byte(updatedConfig), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}
