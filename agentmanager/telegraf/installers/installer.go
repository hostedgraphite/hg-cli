package installers

import (
	"fmt"
	"hg-cli/sysinfo"
	"os"
	"os/exec"
	"strings"
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

func TelegrafPluginInstall(configPath, telegrafCmd string, selectedPlugins []string, sysinfo sysinfo.SysInfo) error {
	plugins := strings.Join(selectedPlugins, ":")
	output := "--output-filter"
	input := "--input-filter"

	commandErr := RunPluginConfig(telegrafCmd, configPath, input, plugins, output, "graphite", "config")
	if commandErr != nil {
		return commandErr
	}

	return nil
}

func TelegrafGraphiteUpdate(apikey, configPath string) error {
	var err error

	currentFilePath := configPath
	newServerURL := "carbon.sandbox.hostedgraphite.com:2003"
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

func RunPluginConfig(command, path string, args ...string) error {

	file, err := os.Create(path)
	if err != nil {

		return fmt.Errorf("error creating output file: %v", err)
	}

	defer file.Close()

	cmd := exec.Command(command, args...)
	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {

		return fmt.Errorf("error executing command: %v", err)
	}

	return nil
}
