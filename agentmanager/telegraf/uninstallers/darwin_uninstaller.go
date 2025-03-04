package uninstallers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf/utils"
)

func DarwinUninstall(pkgMngr, arch string, updates chan<- string) error {
	var err error

	if pkgMngr == "brew" {
		err = brewUninstall(updates)
	} else {
		err = manualUninstall(updates)
	}
	return err
}

func brewUninstall(updates chan<- string) error {
	var err error
	err = utils.RunCommand("brew", []string{"uninstall", "telegraf"}, updates)
	if err != nil {
		return fmt.Errorf("error uninstalling telegraf service: %v", err)
	}

	return err

}

func manualUninstall(updates chan<- string) error {
	var err error

	appPath := "/Applications/Telegraf.app"
	err = os.RemoveAll(appPath)
	if err != nil {
		return fmt.Errorf("error removing Telegraf.app: %w", err)
	}

	telegrafPath := "/usr/local/bin/telegraf"
	err = os.Remove(telegrafPath)
	if err != nil {
		return fmt.Errorf("error removing telegraf binary: %w", err)
	}

	return err
}

func DarwinDeleteFiles(arch string, updates chan<- string) error {
	var err error
	var pattern string

	if arch == "arm64" {
		pattern = "/opt/homebrew/etc/telegraf*"
	} else {
		pattern = "/usr/local/etc/telegraf*"
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("error matching pattern '%s': %v", pattern, err)
	}

	for _, match := range matches {
		err := os.Remove(match)
		if err != nil {
			return fmt.Errorf("error deleting file '%s': %v", match, err)
		}
	}

	return err
}
