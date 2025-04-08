package otel

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	otelPipes "github.com/hostedgraphite/hg-cli/agentmanager/otel/pipes"
	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/pipeline"
)

//go:embed config.yaml
var configYaml []byte

//go:embed com.otelcol-contrib-agent.plist
var plistFile []byte

//go:embed systemd.conf
var systemdFile []byte

func (o *Otel) InstallPipeline(updates chan *pipeline.Pipe) (*pipeline.Pipeline, error) {
	var err error
	var sysInfo = o.sysinfo
	var pipes []*pipeline.Pipe

	switch sysInfo.Os {
	case "linux":
		pipes = otelPipes.LinuxInstallPipes(sysInfo)
	case "darwin":
		pipes = otelPipes.DarwinInstallPipes(sysInfo)
	case "windows":
		pipes, err = otelPipes.WindowsInstallPipes(sysInfo)
		if err != nil {
			return nil, err
		}
	}

	configPipes, err := o.configPipeline()
	if err != nil {
		return nil, err
	}

	pipes = append(pipes, configPipes...)

	pipeline := pipeline.NewPipeline(
		fmt.Sprintf("Installing Otel Agent (%s-%s)",
			sysInfo.Os,
			sysInfo.PkgMngr,
		),
		pipes,
		updates,
	)

	return &pipeline, err
}

func (o *Otel) configPipeline() ([]*pipeline.Pipe, error) {
	var err error
	var pipes []*pipeline.Pipe

	if o.sysinfo.Os == "windows" {
		pipes = otelPipes.WindowsConfigPipes(o.options, o.serviceSettings)
	} else if o.sysinfo.Os == "darwin" {
		pipes = otelPipes.DarwinConfigPipes(o.options, o.serviceSettings, string(plistFile))
	} else if o.sysinfo.Os == "linux" && o.sysinfo.PkgMngr == "" {
		pipes = otelPipes.LinuxManualConfigPipes(o.options, o.serviceSettings, string(systemdFile))
	}

	updatePipe := o.graphiteOutputUpdatePipe()

	pipes = append(pipes, updatePipe...)

	return pipes, err
}

func (o *Otel) graphiteOutputUpdatePipe() []*pipeline.Pipe {
	os := o.sysinfo.Os
	var cmd *exec.Cmd
	var configPath string

	if os == "windows" {
		cmd = exec.Command("powershell", "-Command", "echo test")
	} else {
		cmd = exec.Command("sleep", "1")
	}

	if o.options["configPath"] != nil {
		configPath = o.options["configPath"].(string)
	} else {
		configPath = o.serviceSettings["configPath"]
	}

	pipes := []*pipeline.Pipe{
		pipeline.NewPipe("Updating Otel config.yaml", cmd).PostRun(
			func(ctx context.Context) error {
				return graphiteOutputUpdate(o.apikey, configPath)
			},
		),
	}

	return pipes

}

func graphiteOutputUpdate(apikey, configPath string) error {
	configYaml := string(configYaml)
	graphiteBlock := `(?m)^processors:\n(?:\s{2,}.*\n?)*`

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
	}

	updates := map[string]string{
		`<HG-API-KEY>.*`: fmt.Sprintf(`"%s.opentel."`, apikey), // Works without regex issues
		`<HOSTNAME>`:     hostname,
	}

	updatedConfig, err := utils.UpdateConfigBlock(configYaml, graphiteBlock, updates)

	if err != nil {
		return fmt.Errorf("error during updating: %v", err)
	}

	// $$ causes issued with regex as it's seen as an escape character.
	updatedConfig = strings.ReplaceAll(updatedConfig, "opentel.", "opentel.$$0")

	err = os.WriteFile(configPath, []byte(updatedConfig), 0644)

	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func (o *Otel) UninstallPipeline(updates chan *pipeline.Pipe) (*pipeline.Pipeline, error) {
	var err error
	var sysInfo = o.sysinfo
	var pipes []*pipeline.Pipe

	switch sysInfo.Os {
	case "darwin":
		pipes = otelPipes.DarwinUninstallPipes()
	case "linux":
		pipes = otelPipes.LinxUninstallPipes(sysInfo)
	case "windows":
		pipes, err = otelPipes.WindowsUninstallPipes(sysInfo)
		if err != nil {
			return nil, err
		}
	}

	pipeline := pipeline.NewPipeline(
		fmt.Sprintf("Uninstalling Otel Agent (%s-%s)",
			sysInfo.Os,
			sysInfo.PkgMngr,
		),
		pipes,
		updates,
	)

	return &pipeline, err
}
func (o *Otel) UpdateApiKeyPipeline(updates chan *pipeline.Pipe) (*pipeline.Pipeline, error) {
	var err error
	var sysInfo = o.sysinfo
	var pipes []*pipeline.Pipe

	switch sysInfo.Os {
	case "linux", "darwin", "windows":
		pipes = o.graphiteOutputUpdatePipe()
	default:
		return nil, fmt.Errorf("unsupported operating system: %v", err)
	}

	pipeline := pipeline.NewPipeline(fmt.Sprintf("Updating HostedGraphite Api Key (%s-%s)", sysInfo.Os, sysInfo.PkgMngr), pipes, updates)

	return &pipeline, err

}
