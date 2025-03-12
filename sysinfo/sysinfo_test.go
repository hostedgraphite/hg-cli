package sysinfo

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

var orininalExecCommand = execCommand

const UbunutRelease = `
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

func TestCheckUbuntuDistroPkgMngr(t *testing.T) {
	expectedOs, expectedPkgMngr := "ubuntu", "apt"
	distro, pkgMngr := checkDistroPkgMngr(UbunutRelease)
	require.Equal(t, expectedOs, distro)
	require.Equal(t, expectedPkgMngr, pkgMngr)

}

const FedoraRelease = `
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

func TestFedoraDistorPkg(t *testing.T) {
	expectedOs, expectedPkgMngr := "fedora", "dnf"
	distro, pkgMngr := checkDistroPkgMngr(FedoraRelease)
	require.Equal(t, expectedOs, distro)
	require.Equal(t, expectedPkgMngr, pkgMngr)
}

const CentOSRelease = `
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

func TestRhelDistroPkg(t *testing.T) {
	expectedOs, expectedPkgMngr := "centos", "yum"
	distro, pkgMngr := checkDistroPkgMngr(CentOSRelease)
	require.Equal(t, expectedOs, distro)
	require.Equal(t, expectedPkgMngr, pkgMngr)
}

func TestBadDistroPkg(t *testing.T) {
	distro, pkgMngr := checkDistroPkgMngr("")
	require.Equal(t, "", distro)
	require.Equal(t, "", pkgMngr)
}

func TestCheckPkgMngrTrue(t *testing.T) {
	defer func() { execCommand = orininalExecCommand }()
	output := "/usr/local/bin/brew"
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", output)
	}
	result := checkPkgMngr("brew")
	require.Equal(t, true, result)
}

func TestCheckPkgMngrFalse(t *testing.T) {
	defer func() { execCommand = orininalExecCommand }()
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	result := checkPkgMngr("brew")
	require.Equal(t, false, result)
}
