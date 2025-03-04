package uninstallers

import (
	"fmt"
	"hg-cli/sysinfo"
)

func TelegrafUninstall(sysinfo sysinfo.SysInfo, updates chan<- string) error {
	var err error

	operatingSystem := sysinfo.Os
	arch := sysinfo.Arch
	distro := sysinfo.Distro
	pkgMngr := sysinfo.PkgMngr

	switch sysinfo.Os {
	case "darwin":
		err = DarwinUninstall(pkgMngr, arch, updates)
	case "linux":
		err = LinuxUninstall(operatingSystem, arch, distro, pkgMngr, updates)
	case "windows":
		err = WindowsUninstall(updates)
	default:
		err = fmt.Errorf("unsupported OS: %s", sysinfo.Os)
	}

	return err
}

func TelegrafDeleteFiles(sysinfo sysinfo.SysInfo, updates chan<- string) error {
	var err error
	arch := sysinfo.Arch

	switch sysinfo.Os {
	case "darwin":
		err = DarwinDeleteFiles(arch, updates)
	case "linux":
		err = LinuxDeleteFiles(updates)
	case "windows":
		err = WindowsDeleteFiles(updates)
	default:
		err = fmt.Errorf("unsupported OS: %s", sysinfo.Os)
	}

	return err
}
