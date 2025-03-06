package installers

import (
	"fmt"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
)

var linuxArchFile = map[string]string{
	"amd64": "telegraf-1.33.1_linux_amd64.tar.gz",
	"386":   "telegraf-1.33.1_linux_i386.tar.gz",
	"armv8": "telegraf-1.33.1_linux_arm64.tar.gz",
	"armv7": "telegraf-1.33.1_linux_armhf.tar.gz",
}

func TelegrafAgentInstallLinux(operatingSystem, arch, distro, pkgMngr string, updates chan<- string) error {
	var err error

	if distro == "ubuntu" || distro == "debian" && pkgMngr != "" {
		err = UbuntuDebInstall(updates)
	} else if distro == "redhat" || distro == "centos" || distro == "rhel" && pkgMngr != "" {
		err = CentOsRhelInstall(updates)
	} else {
		err = LinuxInstaller(operatingSystem, arch, updates)
	}

	return err
}

func UbuntuDebInstall(updates chan<- string) error {
	var err error
	if err = utils.RunCommand("sudo", []string{"curl", "--silent", "--location", "-O", "https://repos.influxdata.com/influxdata-archive.key"}, updates); err != nil {
		return fmt.Errorf("failed to download GPG key: %v", err)
	}

	verifyCmd := "echo '943666881a1b8d9b849b74caebf02d3465d6beb716510d86a39f6c8e8dac7515  influxdata-archive.key' | sha256sum -c"
	if err = utils.RunCommand("sudo", []string{"bash", "-c", verifyCmd}, updates); err != nil {
		return fmt.Errorf("failed to verify GPG key: %v", err)
	}

	addKeyCmd := "cat influxdata-archive.key | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/influxdata-archive.gpg > /dev/null"
	if err = utils.RunCommand("sudo", []string{"bash", "-c", addKeyCmd}, updates); err != nil {
		return fmt.Errorf("failed to add GPG key to trusted keys: %v", err)
	}

	addRepoCmd := "echo 'deb [signed-by=/etc/apt/trusted.gpg.d/influxdata-archive.gpg] https://repos.influxdata.com/debian stable main' | sudo tee /etc/apt/sources.list.d/influxdata.list"
	if err = utils.RunCommand("sudo", []string{"bash", "-c", addRepoCmd}, updates); err != nil {
		return fmt.Errorf("failed to add InfluxData repository: %v", err)
	}

	if err = utils.RunCommand("sudo", []string{"apt-get", "update"}, updates); err != nil {
		return fmt.Errorf("failed to update package list: %v", err)
	}

	if err = utils.RunCommand("sudo", []string{"apt-get", "install", "-y", "telegraf"}, updates); err != nil {
		return fmt.Errorf("failed to install Telegraf: %v", err)
	}

	return err
}

func CentOsRhelInstall(updates chan<- string) error {
	var err error
	err = getRepo(updates)

	if err != nil {
		return err
	}

	if err = utils.RunCommand("sudo", []string{"yum", "install", "-y", "telegraf"}, updates); err != nil {
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
			fmt.Sprintf("echo '%s' | sudo tee /etc/yum.repos.d/influxdata.repo", repo),
		},
		updates,
	)

	if err != nil {
		return err
	}

	return err

}

func LinuxInstaller(operatingSystem, arch string, updates chan<- string) error {
	file := linuxArchFile[arch]
	url := "https://dl.influxdata.com/telegraf/releases/" + file

	if err := utils.RunCommand("wget", []string{url}, updates); err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}

	if err := utils.RunCommand("tar", []string{"xf", file}, updates); err != nil {
		return fmt.Errorf("error running tar on file: %v", err)
	}

	if err := utils.RunCommand("mkdir", []string{"/etc/telegraf"}, updates); err != nil {
		return fmt.Errorf("error making dir: %v", err)
	}

	if err := utils.RunCommand("mv", []string{"telegraf-1.33.1/etc/telegraf/telegraf.conf", "/etc/telegraf/"}, updates); err != nil {
		return fmt.Errorf("error moving conf file: %v", err)
	}

	if err := utils.RunCommand("mv", []string{"telegraf-1.33.1/usr/bin/telegraf", "/usr/local/bin"}, updates); err != nil {
		return fmt.Errorf("error moving exe file: %v", err)
	}

	return nil
}
