package uninstallers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf/installers"
)

func WindowsUninstall(updates chan<- string) error {
	if !installers.IsInstalledWindows() {
		updates <- "No exe found - Unable to remove service"
		return nil
	}

	// Stop telegraf windows service
	updates <- "Stopping telegraf service"
	stopservice := `& "C:\Program Files\InfluxData\telegraf\telegraf.exe" --service-name telegraf service stop`
	cmd := exec.Command("powershell", "-Command", stopservice)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "does not exist") {
			updates <- "Not installed as a service"
			return nil
		} else {
			return fmt.Errorf("error stopping telegraf service: %v\nOutput: %s", err, output)
		}
	}

	// Uninstall telegraf windows service
	updates <- "Uninstalling telegraf service"
	uninstall := `& "C:\Program Files\InfluxData\telegraf\telegraf.exe" --service-name telegraf service uninstall`
	cmd = exec.Command("powershell", "-Command", uninstall)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error uninstalling telegraf: %v\nOutput: %s", err, output)
	}
	return err
}

func WindowsDeleteFiles(updates chan<- string) error {
	if !installers.IsInstalledWindows() {
		updates <- "No files found at C:\\Program Files\\InfluxData\\telegraf"
		return nil
	}

	// Delete all files in telegraf directory
	updates <- "Removing telegraf files from C:\\Program Files\\InfluxData\\telegraf"
	telegrafDir := "C:\\Program Files\\InfluxData\\telegraf"
	err := os.RemoveAll(telegrafDir)
	if err != nil {
		return fmt.Errorf("error removing telegraf directory: %w", err)
	}
	return err
}
