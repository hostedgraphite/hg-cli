package sysinfo

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/x/term"
)

type SysInfo struct {
	Os      string
	Arch    string
	PkgMngr string
	Distro  string
	Width   int
	Height  int
}

func checkMacPkgMngr() (string, string) {
	cmd := exec.Command("which", "brew")
	err := cmd.Run()
	if err != nil {
		return "", ""
	}
	return "", "brew"
}

func checkDistroPkgMngr() (string, string) {
	cmd := exec.Command("cat", "/etc/os-release")
	output, err := cmd.Output()
	if err != nil {
		return "", ""
	}

	var distribution string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			distribution = strings.TrimPrefix(line, "ID=")
			distribution = strings.Trim(distribution, `"`) // Remove quotes if present
			break
		}
	}

	if distribution == "" {
		return "", ""
	}

	var packageManager string
	if distribution == "ubuntu" || distribution == "debian" {
		packageManager = "apt"
	} else if distribution == "centos" || distribution == "redhat" || distribution == "rhel" {
		packageManager = "yum"
	} else if distribution == "fedora" {
		packageManager = "dnf"
	} else {
		return "", ""
	}

	cmd = exec.Command("which", packageManager)
	err = cmd.Run()

	if err != nil {
		return distribution, ""
	}

	return distribution, packageManager
}

func GetSystemInformation() (SysInfo, error) {
	var distro, pkgmngr string

	goOs := runtime.GOOS
	goArch := runtime.GOARCH

	if goOs == "darwin" {
		distro, pkgmngr = checkMacPkgMngr()
	} else if goOs == "linux" {
		distro, pkgmngr = checkDistroPkgMngr()
	} else {
		pkgmngr = ""
		distro = ""
	}

	initialHeight, initialWidth := GetInitialDimensions()

	system := SysInfo{
		Os:      strings.ToLower(goOs),
		Arch:    strings.ToLower(goArch),
		PkgMngr: strings.ToLower(pkgmngr),
		Distro:  strings.ToLower(distro),
		Width:   initialWidth,
		Height:  initialHeight,
	}

	return system, nil
}

func GetInitialDimensions() (int, int) {
	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		return 80, 30
	}
	return width, height
}
