package installers

import (
	"fmt"

	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf/utils"
)

func TelegrafAgentInstallDarwin(pkgMngr, arch string, updates chan<- string) error {
	var err error
	if pkgMngr == "" {
		err = MacBinInstall(arch, updates)
	} else {
		err = MacBrewInstall(updates)
	}
	return err
}

func MacBrewInstall(updates chan<- string) error {
	var err error
	if err = utils.RunCommand("brew", []string{"install", "telegraf"}, updates); err != nil {
		return fmt.Errorf("error installing telegraf service: %v", err)
	}
	return err
}

func MacBinInstall(arch string, updates chan<- string) error {
	var dmgURL, dmgFileName string

	// Set the download URL and file name based on architecture
	if arch == "arm64" {
		dmgURL = "https://dl.influxdata.com/telegraf/releases/telegraf-1.33.1_darwin_arm64.dmg"
		dmgFileName = "telegraf-1.33.1_darwin_arm64.dmg"
	} else {
		dmgURL = "https://dl.influxdata.com/telegraf/releases/telegraf-1.33.1_darwin_amd64.dmg"
		dmgFileName = "telegraf-1.33.1_darwin_amd64.dmg"
	}

	volumeName := "/Volumes/Telegraf"

	// Download the DMG
	if err := utils.RunCommand("curl", []string{"-L", dmgURL, "-o", dmgFileName}, updates); err != nil {
		return fmt.Errorf("error downloading Telegraf: %v", err)
	}

	// Attach the DMG
	if err := utils.RunCommand("hdiutil", []string{"attach", dmgFileName}, updates); err != nil {
		return fmt.Errorf("error mounting the DMG file: %v", err)
	}

	// Move the application to /Applications
	if err := utils.RunCommand("cp", []string{"-R", volumeName + "/Telegraf.app", "/Applications/"}, updates); err != nil {
		return fmt.Errorf("error moving Telegraf application: %v", err)
	}

	// Copy the binary to /usr/local/bin
	if err := utils.RunCommand("cp", []string{volumeName + "/Telegraf.app/Contents/Resources/usr/bin/telegraf", "/usr/local/bin/"}, updates); err != nil {
		return fmt.Errorf("error copying telegraf binary: %v", err)
	}

	// Detach the DMG
	if err := utils.RunCommand("hdiutil", []string{"detach", volumeName}, updates); err != nil {
		return fmt.Errorf("error unmounting the DMG file: %v", err)
	}

	return nil
}
