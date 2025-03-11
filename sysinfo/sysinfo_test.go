package sysinfo

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

var orininalExecCommand = execCommand
var originalOsGeteUid = osGeteUid

func TestCheckUbuntuDistroPkgMngr(t *testing.T) {
	defer func() { execCommand = orininalExecCommand }()
	output := `
NAME="Ubuntu"
VERSION="20.04.5 LTS (Focal Fossa)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 20.04.5 LTS"
VERSION_ID="20.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=focal
UBUNTU_CODENAME=focal
`

	expectedOs, expectedPkgMngr := "ubuntu", "apt"

	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", output)
	}

	os, pkgMngr := checkDistroPkgMngr()
	require.Equal(t, expectedOs, os)
	require.Equal(t, expectedPkgMngr, pkgMngr)

}

func TestFedoraDistorPkg(t *testing.T) {
	defer func() { execCommand = orininalExecCommand }()
	output := `
NAME=Fedora
VERSION="28 (Twenty Eight)"
ID=fedora
VERSION_ID=28
VERSION_CODENAME=""
PLATFORM_ID="platform:f28"
PRETTY_NAME="Fedora 28 (Twenty Eight)"
ANSI_COLOR="0;34"
LOGO=fedora-logo-icon
CPE_NAME="cpe:/o:fedoraproject:fedora:28"
HOME_URL="https://fedoraproject.org/"
SUPPORT_URL="https://fedoraproject.org/wiki/Communicating_and_getting_help"
BUG_REPORT_URL="https://bugzilla.redhat.com/"
REDHAT_BUGZILLA_PRODUCT="Fedora"
REDHAT_BUGZILLA_PRODUCT_VERSION=28
REDHAT_SUPPORT_PRODUCT="Fedora"
REDHAT_SUPPORT_PRODUCT_VERSION=28
PRIVACY_POLICY_URL="https://fedoraproject.org/wiki/Legal:PrivacyPolicy"
`

	expectedOs, expectedPkgMngr := "fedora", "dnf"
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", output)
	}

	os, pkgMngr := checkDistroPkgMngr()
	require.Equal(t, expectedOs, os)
	require.Equal(t, expectedPkgMngr, pkgMngr)
}

func TestRhelDistroPkg(t *testing.T) {
	defer func() { execCommand = orininalExecCommand }()
	output := `
NAME="CentOS Linux"
VERSION="8"
ID="centos"
ID_LIKE="rhel fedora"
VERSION_ID="8"
PLATFORM_ID="platform:el8"
PRETTY_NAME="CentOS Linux 8"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:centos:centos:8"
HOME_URL="https://centos.org/"
BUG_REPORT_URL="https://bugs.centos.org/"
CENTOS_MANTISBT_PROJECT="CentOS-8"
CENTOS_MANTISBT_PROJECT_VERSION="8"
`
	expectedOs, expectedPkgMngr := "centos", "yum"
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", output)
	}
	os, pkgMngr := checkDistroPkgMngr()
	require.Equal(t, expectedOs, os)
	require.Equal(t, expectedPkgMngr, pkgMngr)
}

func TestErrorDistroPkg(t *testing.T) {
	defer func() { execCommand = orininalExecCommand }()
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	os, pkgMngr := checkDistroPkgMngr()
	require.Equal(t, "", os)
	require.Equal(t, "", pkgMngr)
}

func TestMacPkgMngr(t *testing.T) {
	defer func() { execCommand = orininalExecCommand }()
	output := "/usr/local/bin/brew"
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", output)
	}
	_, pkgMngr := checkMacPkgMngr()
	require.Equal(t, "brew", pkgMngr)
}

func TestNoMacPkgMngr(t *testing.T) {
	defer func() { execCommand = orininalExecCommand }()
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	_, pkgMngr := checkMacPkgMngr()
	require.Equal(t, "", pkgMngr)
}

func TestSudoPerm(t *testing.T) {
	defer func() { osGeteUid = originalOsGeteUid }()
	osGeteUid = 0 // sudo permissions
	expected := true

	require.Equal(t, expected, checkSudoPerm())
}

func TestNoSudoPerm(t *testing.T) {
	defer func() { osGeteUid = originalOsGeteUid }()
	osGeteUid = 1000 // no sudo permissions
	expected := false

	require.Equal(t, expected, checkSudoPerm())
}
