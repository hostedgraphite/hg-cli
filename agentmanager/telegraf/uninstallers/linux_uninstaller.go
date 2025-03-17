package uninstallers

import (
	"fmt"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
)

func LinuxUninstall(operatingSystem, arch, distro, pkgMngr string, updates chan<- string) error {
	var err error

	if pkgMngr == "brew" {
		err = BrewUninstall(updates)
	} else if pkgMngr == "" || (distro != "ubuntu" && distro != "debian" && distro != "redhat" && distro != "centos" && distro != "rhel" && distro != "fedora") {
		err = LinuxUninstaller(updates)
	} else {
		err = LinuxPkgMngrUninstaller(pkgMngr, updates)
	}
	return err
}

func LinuxPkgMngrUninstaller(pkgMngr string, updates chan<- string) error {
	if err := utils.RunCommand(pkgMngr, []string{"remove", "telegraf", "-y"}, updates); err != nil {
		return fmt.Errorf("error uninstalling telegraf service: %v", err)
	}
	return nil
}

func LinuxUninstaller(updates chan<- string) error {
	var err error
	telegrafBin := "/usr/bin/telegraf"
	telegrafDir := "/etc/telegraf"
	telegrafSystemd := "/etc/systemd/system/telegraf.service"

	if err = utils.RunCommand("systemctl", []string{"stop", "telegraf"}, updates); err != nil {
		return fmt.Errorf("error stopping telegraf service: %v", err)
	}

	if err = utils.RunCommand("rm", []string{"-rf", telegrafBin}, updates); err != nil {
		return fmt.Errorf("error stopping telegraf service: %v", err)
	}

	if err = utils.RunCommand("rm", []string{"-rf", telegrafDir}, updates); err != nil {
		return fmt.Errorf("error stopping telegraf service: %v", err)
	}

	if err = utils.RunCommand("rm", []string{"-rf", telegrafSystemd}, updates); err != nil {
		return fmt.Errorf("error stopping telegraf service: %v", err)
	}

	if err = utils.RunCommand("userdel", []string{"telegraf"}, updates); err != nil {
		return fmt.Errorf("error stopping telegraf service: %v", err)
	}

	return err
}

func LinuxDeleteFiles(updates chan<- string) error {
	var err error

	telegrafDir := "/etc/telegraf"

	if err = utils.RunCommand("rm", []string{"-rf", telegrafDir}, updates); err != nil {
		return fmt.Errorf("error deleting telegraf files: %v", err)
	}

	return err
}
