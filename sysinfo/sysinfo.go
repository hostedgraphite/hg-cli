package sysinfo

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/x/term"
)

type SysInfo struct {
	Os       string
	Arch     string
	PkgMngr  string
	Distro   string
	SudoPerm bool
	Width    int
	Height   int
}

var execCommand = exec.Command

func checkSudoPerm() bool {
	return os.Getegid() == 0
}

func checkHgCliBrewInstall() bool {
	cmd := execCommand("brew", "list", "--formula", "hg-cli")
	err := cmd.Run()
	return err == nil
}

func checkSudoPermWindows() bool {
	cmd := execCommand("net", "session")
	err := cmd.Run()
	return err == nil
}

func checkPkgMngr(packageManager string) bool {
	cmd := execCommand("which", packageManager)
	err := cmd.Run()
	return err == nil
}

func getOSRelease() (string, error) {
	cmd := execCommand("cat", "/etc/os-release")
	output, err := cmd.Output()
	return string(output), err
}

func checkDistroPkgMngr(releaseInfo string) (string, string) {
	var distribution string
	lines := strings.Split(releaseInfo, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			distribution = strings.TrimPrefix(line, "ID=")
			distribution = strings.Trim(distribution, `"`) // Remove quotes if present
			break
		}
	}

	var packageManager string
	var distroMap = map[string]string{
		"ubuntu": "apt",
		"debian": "apt",
		"redhat": "yum",
		"centos": "yum",
		"rhel":   "yum",
		"fedora": "dnf",
	}
	packageManager = distroMap[distribution]
	return distribution, packageManager
}

func GetSystemInformation() (SysInfo, error) {
	var distro, pkgmngr string
	var sudoPerm bool

	goOs := runtime.GOOS
	goArch := runtime.GOARCH

	// Determine the package manager & distro
	switch goOs {
	case "darwin":
		pkgmngr = "brew"
		sudoPerm = checkSudoPerm()
	case "linux":
		releaseInfo, err := getOSRelease()
		if err == nil {
			distro, pkgmngr = checkDistroPkgMngr(releaseInfo)
			sudoPerm = checkSudoPerm()
		}
	case "windows":
		sudoPerm = checkSudoPermWindows()
	}

	// If hg-cli was installed with brew, set the package manager to brew
	if checkHgCliBrewInstall() {
		pkgmngr = "brew"
	}

	// Confirm that the package manager is installed
	if !checkPkgMngr(pkgmngr) {
		pkgmngr = ""
	}

	initialHeight, initialWidth := GetInitialDimensions()

	system := SysInfo{
		Os:       strings.ToLower(goOs),
		Arch:     strings.ToLower(goArch),
		PkgMngr:  strings.ToLower(pkgmngr),
		Distro:   strings.ToLower(distro),
		SudoPerm: sudoPerm,
		Width:    initialWidth,
		Height:   initialHeight,
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
