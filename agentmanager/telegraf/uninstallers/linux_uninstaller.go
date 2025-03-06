package uninstallers

import (
	"fmt"
	"os"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
)

func LinuxUninstall(operatingSystem, arch, distro, pkgMngr string, updates chan<- string) error {
	var err error

	if distro == "ubuntu" || distro == "debian" && pkgMngr != "" {
		err = UbuntuDebUninstaller(updates)
	} else if distro == "redhat" || distro == "centos" || distro == "rhel" && pkgMngr != "" {
		err = CentOsRhelUninstaller(updates)
	} else {
		err = LinuxUninstaller(updates)
	}

	return err
}

func UbuntuDebUninstaller(updates chan<- string) error {
	var err error

	cmd := utils.RunCommand("sudo", []string{"apt-get", "remove", "telegraf", "-y"}, updates)
	if cmd != nil {
		return fmt.Errorf("error uninstalling telegraf service: %v", cmd)
	}

	return err
}

func CentOsRhelUninstaller(updates chan<- string) error {
	var err error

	cmd := utils.RunCommand("sudo", []string{"yum", "remove", "telegraf", "-y"}, updates)
	if cmd != nil {
		return fmt.Errorf("error uninstalling telegraf service: %v", cmd)
	}

	return err
}

func LinuxUninstaller(updates chan<- string) error {
	var err error
	telegrafBin := "/usr/local/bin/telegraf"
	err = os.Remove(telegrafBin)
	if err != nil {
		return fmt.Errorf("error removing telegraf binary: %w", err)
	}

	return err
}

func LinuxDeleteFiles(updates chan<- string) error {
	var err error

	telegrafDir := "/etc/telegraf"

	if err = utils.RunCommand("sudo", []string{"rm", "-rf", telegrafDir}, updates); err != nil {
		return fmt.Errorf("error deleting telegraf files: %v", err)
	}

	return err
}
