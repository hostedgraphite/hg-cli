package uninstallers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf/installers"
)

func WindowsUninstall(updates chan<- string) error {
	if !installers.IsInstalledWindows() {
		updates <- "Not exe found - Unable to remove service"
		return nil
	}

	uninstall := `& "C:\Program Files\InfluxData\telegraf\telegraf.exe" --service-name telegraf service uninstall`
	cmd := exec.Command("powershell", "-Command", uninstall)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error uninstalling telegraf: %v\nOutput: %s", err, output)
	}
	return err
}

func WindowsDeleteFiles(updates chan<- string) error {
	if !installers.IsInstalledWindows() {
		updates <- "No telegraf files found at C:\\Program Files\\InfluxData\\telegraf"
		return nil
	}

	telegrafDir := "C:\\Program Files\\InfluxData\\telegraf"
	err := os.RemoveAll(telegrafDir)
	if err != nil {
		return fmt.Errorf("error removing telegraf directory: %w", err)
	}
	return err
}
