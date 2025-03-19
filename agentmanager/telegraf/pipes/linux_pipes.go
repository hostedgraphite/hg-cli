package pipes

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/pipeline"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

func LinuxInstallPipes(sysInfo sysinfo.SysInfo) []*pipeline.Pipe {
	var pipes []*pipeline.Pipe
	arch := sysInfo.Arch
	distro := sysInfo.Distro
	pkgMngr := sysInfo.PkgMngr

	if pkgMngr == "brew" {
		pipes = BrewInstallPipes()
	} else if pkgMngr == "apt" {
		pipes = aptInstallPipes()
	} else if pkgMngr == "yum" || pkgMngr == "dnf" {
		pipes = yumInstallPipes()
	} else {
		pipes = linuxBinInstallPipes(arch, distro)
	}

	return pipes
}

func LinuxConfigPipes(options map[string]interface{}, serviceSettings map[string]string) []*pipeline.Pipe {
	inputs := strings.Join(options["plugins"].([]string), ":")
	telegrafCmd := serviceSettings["serviceCmd"]
	configpath := serviceSettings["configPath"]

	pipes := []*pipeline.Pipe{
		{
			Name: "Configuring Telegraf Plugins",
			Cmd:  exec.Command("sh", "-c", telegrafCmd+" --input-filter "+inputs+" --output-filter graphite config > "+configpath),
		},
	}
	return pipes
}

func aptInstallPipes() []*pipeline.Pipe {

	tmpDir := "/tmp/hg-cli"
	keyPath := "/tmp/hg-cli/influxdata-archive.key"

	pipes := []*pipeline.Pipe{
		{
			Name: "Creating TMP Directory",
			Cmd:  exec.Command("mkdir", "-p", tmpDir),
		},
		{
			Name: "Getting Influx archive Key",
			Cmd:  exec.Command("curl", "--silent", "--location", "-o", keyPath, "https://repos.influxdata.com/influxdata-archive.key"),
		},
		{
			Name: "Adding Influx archive Key to apt trusted",
			Cmd:  exec.Command("bash", "-c", fmt.Sprintf("cat %s | gpg --dearmor > /etc/apt/trusted.gpg.d/influxdata-archive.gpg", keyPath)),
		},
		{
			Name: "Adding InfluxData apt Repository",
			Cmd:  exec.Command("bash", "-c", "echo 'deb [signed-by=/etc/apt/trusted.gpg.d/influxdata-archive.gpg] https://repos.influxdata.com/debian stable main' > /etc/apt/sources.list.d/influxdata.list"),
		},
		{
			Name: "Updating Package List",
			Cmd:  exec.Command("apt-get", "update"),
		},
		{
			Name: "Installing Telegraf",
			Cmd:  exec.Command("apt-get", "install", "-y", "telegraf"),
		},
		{
			Name: "Deleting TMP Directory",
			Cmd:  exec.Command("rm", "-rf", tmpDir),
		},
	}

	return pipes
}

const yumRepo = `[influxdata]
name = InfluxData Repository - Stable
baseurl = https://repos.influxdata.com/stable/$basearch/main
enabled = 1
gpgcheck = 1
gpgkey = https://repos.influxdata.com/influxdata-archive_compat.key`

func yumInstallPipes() []*pipeline.Pipe {

	pipes := []*pipeline.Pipe{
		{
			Name: "Adding InfluxData yum Repository",
			Cmd:  exec.Command("sh", "-c", "echo '"+yumRepo+"' > /etc/yum.repos.d/influxdata.repo"),
		},
		{
			Name: "Installing Telegraf Agent",
			Cmd:  exec.Command("yum", "install", "-y", "telegraf"),
		},
	}

	return pipes
}

var linuxArchFile = map[string]string{
	"amd64":  "_linux_amd64.tar.gz",
	"386":    "_linux_i386.tar.gz",
	"arm64":  "_linux_arm64.tar.gz",
	"armv7l": "_linux_armhf.tar.gz",
}

func linuxBinInstallPipes(arch, distro string) []*pipeline.Pipe {
	var pipes []*pipeline.Pipe

	latest, err := utils.GetLatestReleaseTag("influxdata", "telegraf")
	if err != nil {
		latest = "v1.33.1" // Default
	}
	latest = latest[1:]

	file := "telegraf-" + latest + linuxArchFile[arch]
	url := "https://dl.influxdata.com/telegraf/releases/" + file
	version := "telegraf-" + latest + "/"
	tmpDir := "/tmp/hg-cli/"
	tmpPath := "/tmp/hg-cli/" + file
	telegrafPath := tmpDir + version
	telegrafConf := telegrafPath + "etc/telegraf/telegraf.conf"
	telegrafBin := telegrafPath + "usr/bin/telegraf"
	telegrafService := telegrafPath + "usr/lib/telegraf/scripts/telegraf.service"

	pipes = []*pipeline.Pipe{
		{
			Name: "Creating TMP Directory",
			Cmd:  exec.Command("mkdir", "-p", tmpDir),
		},
		{
			Name: "Downloading Telegraf archive file",
			Cmd:  exec.Command("wget", url, "-q", "-O", tmpPath),
		},
		{
			Name: "Extracting Telegraf archive file",
			Cmd:  exec.Command("tar", "xf", tmpPath, "-C", tmpDir),
		},
		{
			Name: "Creating Telegraf Config Directory",
			Cmd:  exec.Command("mkdir", "-p", "/etc/telegraf"),
		},
		{
			Name: "Moving Telegraf Conf File",
			Cmd:  exec.Command("mv", telegrafConf, "/etc/telegraf/"),
		},
		{
			Name: "Placing bin file in /usr/bin",
			Cmd:  exec.Command("mv", telegrafBin, "/usr/bin/"),
		},
		{
			Name: "Adding service file to systemd",
			Cmd:  exec.Command("mv", telegrafService, "/etc/systemd/system/telegraf.service"),
		},
		{
			Name: "Creating telegraf service group",
			Cmd:  exec.Command("groupadd", "-g", "988", "telegraf"),
		},
		{
			Name: "Creating telegraf user",
			Cmd:  exec.Command("useradd", "-r", "-u", "989", "-g", "988", "-d", "/etc/telegraf", "-s", "/bin/false", "telegraf"),
		},
	}

	if distro == "fedora" || distro == "centos" || distro == "rhel" {
		// For Fedora/CentOs SELinux permissions
		pipes = append(pipes, []*pipeline.Pipe{
			{
				Name: "Setting SELinux permissions",
				Cmd:  exec.Command("restorecon", "-Rv", "/usr/bin/telegraf"),
			},
			{
				Name: "Setting SELinux permissions",
				Cmd:  exec.Command("restorecon", "-Rv", "/etc/systemd/system/telegraf.service"),
			},
		}...)
	}

	pipes = append(pipes, []*pipeline.Pipe{
		{
			Name: "Cleaning up temp dir",
			Cmd:  exec.Command("rm", "-rf", tmpDir),
		},
	}...)

	return pipes

}

func LinuxUninstallPipes(sysInfo sysinfo.SysInfo) []*pipeline.Pipe {
	var pipes []*pipeline.Pipe
	pkgMngr := sysInfo.PkgMngr

	if pkgMngr == "brew" {
		pipes = BrewUninstallPipes()
	} else if pkgMngr == "" {
		pipes = linuxUninstallerPipes()
	} else {
		pipes = linuxPkgMngrUninstallPipes(pkgMngr)
	}

	return pipes
}

func linuxPkgMngrUninstallPipes(pkgMngr string) []*pipeline.Pipe {
	pipes := []*pipeline.Pipe{
		{
			Name: "Stopping Telegraf Service",
			Cmd:  exec.Command(pkgMngr, "stop", "telegraf"),
		},
		{
			Name: "Uninstalling Telegraf Agent",
			Cmd:  exec.Command(pkgMngr, "remove", "telegraf", "-y"),
		},
	}

	return pipes
}

func linuxUninstallerPipes() []*pipeline.Pipe {
	pipes := []*pipeline.Pipe{
		{
			Name: "Stopping Telegraf Service",
			Cmd:  exec.Command("systemctl", "stop", "telegraf"),
		},
		{
			Name: "Removing Telegraf Binary",
			Cmd:  exec.Command("rm", "/usr/bin/telegraf"),
		},
		{
			Name: "Removing Telegraf Service",
			Cmd:  exec.Command("rm", "-rf", "/etc/systemd/system/telegraf.service"),
		},
		{
			Name: "Removing Telegraf User",
			Cmd:  exec.Command("userdel", "telegraf"),
		},
	}

	return pipes
}
