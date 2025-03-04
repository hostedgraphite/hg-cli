package uninstallers

import (
	"fmt"
	"hg-cli/agentmanager/telegraf/utils"
	"os"
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

	cmd := utils.RunCommand("apt-get", []string{"remove", "telegraf", "-y"}, updates)
	if cmd != nil {
		return fmt.Errorf("error uninstalling telegraf service: %v", cmd)
	}

	return err
}

func CentOsRhelUninstaller(updates chan<- string) error {
	var err error

	cmd := utils.RunCommand("yum", []string{"remove", "telegraf", "-y"}, updates)
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

	telelegrafDir := "/etc/telegraf"
	err = os.RemoveAll(telelegrafDir)
	if err != nil {
		return fmt.Errorf("error removing telegraf directory: %w", err)
	}

	return err
}
