package installers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func TelegrafAgentInstallWindows(updates chan<- string) error {
	var err error
	// Will throw error if files exists or if telegraf is installed
	fileExists, programInstalled := windowsCheck()
	var commands []string
	if !fileExists && !programInstalled {
		fmt.Println("Downloading telegraf")
		pshellDownload := `
		wget https://dl.influxdata.com/telegraf/releases/telegraf-1.33.0_windows_amd64.zip -OutFile telegraf-1.33.0_windows_amd64.zip;
		`

		pshellExpand := `
		Expand-Archive .\telegraf-1.33.0_windows_amd64.zip -DestinationPath 'C:\Program Files\InfluxData\telegraf\'
		`
		commands = append(commands, pshellDownload, pshellExpand)

	} else if !fileExists {
		pshellMvFiles := `
		Move-Item "C:\Program Files\InfluxData\telegraf\telegraf.*" "C:\Program Files\InfluxData\telegraf\"
		`
		commands = append(commands, pshellMvFiles)
	} else if programInstalled {
		telUninstall := `
		& "C:\Program Files\InfluxData\telegraf\telegraf.exe" --service uninstall;
		`
		commands = append(commands, telUninstall)
	}

	pshellInstall := `
	& "C:\Program Files\InfluxData\telegraf\telegraf.exe" --service install --config "C:\Program Files\InfluxData\telegraf\telegraf.conf"
	`

	commands = append(commands, pshellInstall)
	for _, commcommands := range commands {
		cmd := exec.Command("powershell", "-Command", commcommands)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error installing telegraf: %v\nOutput: %s", err, output)
		}
	}

	return err
}

func checkFileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func checkProgramInstalled(program string) bool {
	cmd := exec.Command("powershell", "-Command", "Get-Command", program, "-ErrorAction SilentlyContinue")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false
	}
	return strings.Contains(string(output), program)
}

func windowsCheck() (bool, bool) {
	var fileExists, programInstalled bool
	// Check if specific files exist
	filesToCheck := []string{
		"C:\\Program Files\\InfluxData\\telegraf\\telegraf.exe",
		"C:\\Program Files\\InfluxData\\telegraf\\telegraf.conf",
	}

	for _, file := range filesToCheck {
		if checkFileExists(file) {
			fileExists = true
		}
	}

	// Check if the program is installed
	if checkProgramInstalled("telegraf") {
		programInstalled = true
	}

	return fileExists, programInstalled
}
