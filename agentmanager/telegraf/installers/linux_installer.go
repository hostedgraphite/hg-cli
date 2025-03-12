package installers

import (
	"fmt"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
)

var linuxArchFile = map[string]string{
	"amd64":  "telegraf-1.33.1_linux_amd64.tar.gz",
	"386":    "telegraf-1.33.1_linux_i386.tar.gz",
	"arm64":  "telegraf-1.33.1_linux_arm64.tar.gz",
	"armv7l": "telegraf-1.33.1_linux_armhf.tar.gz",
}

func TelegrafAgentInstallLinux(operatingSystem, arch, distro, pkgMngr string, updates chan<- string) error {
	var err error

	if distro == "ubuntu" || distro == "debian" && pkgMngr != "" {
		err = UbuntuDebInstall(updates)
	} else if (distro == "redhat" || distro == "centos" || distro == "rhel" || distro == "fedora") && pkgMngr != "" {
		err = CentOsRhelInstall(updates)
	} else {
		err = LinuxInstaller(operatingSystem, arch, distro, updates)
	}

	return err
}

func UbuntuDebInstall(updates chan<- string) error {
	var err error

	tmpDir := "/tmp/hg-cli"
	keyPath := "/tmp/hg-cli/influxdata-archive.key"

	if err = utils.RunCommand("mkdir", []string{"-p", tmpDir}, updates); err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}

	if err = utils.RunCommand("curl", []string{"--silent", "--location", "-o", keyPath, "https://repos.influxdata.com/influxdata-archive.key"}, updates); err != nil {
		return fmt.Errorf("failed to download GPG key: %v", err)
	}

	addKeyCmd := fmt.Sprintf("cat %s | gpg --dearmor | tee /etc/apt/trusted.gpg.d/influxdata-archive.gpg > /dev/null", keyPath)
	if err = utils.RunCommand("bash", []string{"-c", addKeyCmd}, updates); err != nil {
		return fmt.Errorf("failed to add GPG key to trusted keys: %v", err)
	}

	addRepoCmd := "echo 'deb [signed-by=/etc/apt/trusted.gpg.d/influxdata-archive.gpg] https://repos.influxdata.com/debian stable main' | tee /etc/apt/sources.list.d/influxdata.list"
	if err = utils.RunCommand("bash", []string{"-c", addRepoCmd}, updates); err != nil {
		return fmt.Errorf("failed to add InfluxData repository: %v", err)
	}

	if err = utils.RunCommand("apt-get", []string{"update"}, updates); err != nil {
		return fmt.Errorf("failed to update package list: %v", err)
	}

	if err = utils.RunCommand("apt-get", []string{"install", "-y", "telegraf"}, updates); err != nil {
		return fmt.Errorf("failed to install Telegraf: %v", err)
	}

	if err = utils.RunCommand("rm", []string{"-rf", tmpDir}, updates); err != nil {
		return fmt.Errorf("failed to delete tmp dir %v", err)
	}

	return err
}

func CentOsRhelInstall(updates chan<- string) error {
	var err error
	err = getRepo(updates)

	if err != nil {
		return err
	}

	if err = utils.RunCommand("yum", []string{"install", "-y", "telegraf"}, updates); err != nil {
		return fmt.Errorf("error installing telegraf: %v", err)
	}

	return err
}

func getRepo(updates chan<- string) error {
	var err error
	repo := `[influxdata]
name = InfluxData Repository - Stable
baseurl = https://repos.influxdata.com/stable/$basearch/main
enabled = 1
gpgcheck = 1
gpgkey = https://repos.influxdata.com/influxdata-archive_compat.key`

	err = utils.RunCommand(
		"sh",
		[]string{
			"-c",
			fmt.Sprintf("echo '%s' | tee /etc/yum.repos.d/influxdata.repo", repo),
		},
		updates,
	)

	if err != nil {
		return err
	}

	return err

}

func LinuxInstaller(operatingSystem, arch, distro string, updates chan<- string) error {
	file := linuxArchFile[arch]
	url := "https://dl.influxdata.com/telegraf/releases/" + file
	version := "telegraf-1.33.1/"
	tmpDir := "/tmp/hg-cli/"
	tmpPath := "/tmp/hg-cli/" + file

	if err := utils.RunCommand("mkdir", []string{tmpDir}, updates); err != nil {
		return fmt.Errorf("error creating temp dir")
	}

	if err := utils.RunCommand("wget", []string{url, "-q", "-O", tmpPath}, updates); err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}

	if err := utils.RunCommand("tar", []string{"xf", tmpPath, "-C", tmpDir}, updates); err != nil {
		return fmt.Errorf("error running tar on file: %v", err)
	}

	if err := utils.RunCommand("mkdir", []string{"/etc/telegraf"}, updates); err != nil {
		return fmt.Errorf("error making dir: %v", err)
	}

	telegrafPath := tmpDir + version
	telegrafConf := telegrafPath + "etc/telegraf/telegraf.conf"
	telegrafBin := telegrafPath + "usr/bin/telegraf"
	telegrafService := telegrafPath + "usr/lib/telegraf/scripts/telegraf.service"

	if err := utils.RunCommand("mv", []string{telegrafConf, "/etc/telegraf/"}, updates); err != nil {
		return fmt.Errorf("error moving conf file: %v", err)
	}

	if err := utils.RunCommand("mv", []string{telegrafBin, "/usr/bin/"}, updates); err != nil {
		return fmt.Errorf("error moving exe file: %v", err)
	}

	if err := utils.RunCommand("mv", []string{telegrafService, "/etc/systemd/system/telegraf.service"}, updates); err != nil {
		return fmt.Errorf("error moving service file: %v", err)
	}

	// create telegraf user
	if err := utils.RunCommand("groupadd", []string{"-g", "988", "telegraf"}, updates); err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}
	if err := utils.RunCommand("useradd", []string{"-r", "-u", "989", "-g", "988", "-d", "/etc/telegraf", "-s", "/bin/false", "telegraf"}, updates); err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	if distro == "fedora" || distro == "centos" || distro == "rhel" {
		// For Fedora/CentOs SELinux permissions
		if err := utils.RunCommand("restorecon", []string{"-Rv", "/usr/bin/telegraf"}, updates); err != nil {
			return fmt.Errorf("error setting SELinux permissions: %v", err)
		}
		if err := utils.RunCommand("restorecon", []string{"-Rv", "/etc/systemd/system/telegraf.service"}, updates); err != nil {
			return fmt.Errorf("error setting SELinux permissions: %v", err)
		}
	}

	if err := utils.RunCommand("rm", []string{"-rf", tmpDir}, updates); err != nil {
		return fmt.Errorf("error cleaning up temp dir")
	}

	return nil
}
