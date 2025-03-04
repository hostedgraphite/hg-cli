package uninstallers

import (
	"fmt"
	"os"
	"os/exec"
)

func WindowsUninstall(updates chan<- string) error {
	uninstall := `& "C:\Program Files\InfluxData\telegraf\telegraf.exe" --service uninstall`
	cmd := exec.Command("powershell", "-Command", uninstall)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error uninstalling telegraf: %v\nOutput: %s", err, output)
	}
	return err
}

func WindowsDeleteFiles(updates chan<- string) error {
	telegrafDir := "C:\\Program Files\\InfluxData"
	err := os.RemoveAll(telegrafDir)
	if err != nil {
		return fmt.Errorf("error removing telegraf directory: %w", err)
	}
	return err
}
