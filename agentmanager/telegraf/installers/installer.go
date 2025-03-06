package installers

import (
	"fmt"
	"os"

	"strings"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

func TelegrafAgentInstall(sysinfo sysinfo.SysInfo, updates chan<- string) error {
	var err error
	operatingSystem := sysinfo.Os
	arch := sysinfo.Arch
	distro := sysinfo.Distro
	pkgMngr := sysinfo.PkgMngr

	switch operatingSystem {
	case "darwin":
		err = TelegrafAgentInstallDarwin(pkgMngr, arch, updates)
	case "linux":
		err = TelegrafAgentInstallLinux(operatingSystem, arch, distro, pkgMngr, updates)
	case "windows":
		err = TelegrafAgentInstallWindows(updates)
	default:
		err = fmt.Errorf("unsupported operating system: %v", err)
	}

	return err
}

func TelegrafPluginInstall(configPath, telCmd string, selectedPlugins []string, sysinfo sysinfo.SysInfo, updates chan<- string) error {
	plugins := strings.Join(selectedPlugins, ":")
	output := "--output-filter"
	input := "--input-filter"

	commandErr := RunPluginConfig(telCmd, input, plugins, output, configPath, sysinfo.Os, updates)
	if commandErr != nil {
		return commandErr
	}

	return nil
}

func TelegrafGraphiteUpdate(apikey, configPath, opersystem string, updates chan<- string) error {
	var err error

	if opersystem == "linux" {
		err = linuxGraphiteUpdate(apikey, configPath)
		if err != nil {
			return err
		}
		return nil
	}

	err = graphiteUpdater(apikey, configPath)
	if err != nil {
		return err
	}

	return nil
}

func linuxGraphiteUpdate(apikey, configPath string) error {
	fullConfig, err := utils.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	graphiteBlock := `\[\[outputs\.graphite\]\](?:.|\s)*?\[\[`

	updates := map[string]string{
		`prefix\s*=\s*".*?"`:    fmt.Sprintf(`prefix = "%s.telegraf"`, apikey),
		`servers\s*=\s*\[.*?\]`: `servers = ["carbon.hostedgraphite.com:2003"]`,
		`template\s*=\s*".*?"`:  `## template = "host.tags.measurement.field"`,
	}
	updatedConfig, err := utils.UpdateConfigBlock(fullConfig, graphiteBlock, updates)

	if err != nil {
		return fmt.Errorf("error during updating: %v", err)
	}

	err = utils.WriteFile(configPath, updatedConfig)

	if err != nil {
		return fmt.Errorf("error writing file:%v", err)
	}

	return nil
}

func graphiteUpdater(apikey, configPath string) error {
	var err error

	currentFilePath := configPath
	newServerURL := "carbon.hostedgraphite.com:2003"
	newPrefix := apikey + ".telegraf"
	newTemplate := "## template = host.tags.measurement.field"

	data, err := os.ReadFile(currentFilePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	configStr := string(data)

	configStr = strings.Replace(configStr, `servers = ["localhost:2003"]`, fmt.Sprintf(`servers  = ["%s"]`, newServerURL), 1)
	configStr = strings.Replace(configStr, `prefix = ""`, fmt.Sprintf(`prefix = "%s"`, newPrefix), 1)
	configStr = strings.Replace(configStr, `template = "host.tags.measurement.field"`, newTemplate, 1)

	err = os.WriteFile(currentFilePath, []byte(configStr), 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil

}

func RunPluginConfig(telCmd, input, plugins, output, path, opersystem string, updates chan<- string) error {
	var err error
	var cmd string

	_, err = os.Stat(path)
	if err != nil {
		utils.RunCommand("sudo", []string{"touch", path}, updates)
		utils.RunCommand("sudo", []string{"chmod", "0644", path}, updates)
	}

	if opersystem == "linux" {
		sudoCmd := fmt.Sprintf("sudo tee %s > /dev/null", path)
		cmd = fmt.Sprintf("%s %s %s %s graphite config | %s", telCmd, input, plugins, output, sudoCmd)
		err = utils.RunCommand("sh", []string{"-c", cmd}, updates)
	} else {
		err = utils.RunCommand(telCmd, []string{input, plugins, output, "graphite", "config", ">", path}, updates)
	}

	if err != nil {
		return err
	}

	return err
}
