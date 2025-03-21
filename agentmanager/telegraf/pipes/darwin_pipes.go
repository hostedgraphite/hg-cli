package pipes

import (
	"os/exec"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/pipeline"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

func DarwinInstallPipes(sysInfo sysinfo.SysInfo) []*pipeline.Pipe {
	var pipes []*pipeline.Pipe
	pkgMngr := sysInfo.PkgMngr
	arch := sysInfo.Arch

	if pkgMngr == "brew" {
		pipes = BrewInstallPipes()
	} else {
		pipes = macDmgInstallPipes(arch)
	}

	return pipes
}

func macDmgInstallPipes(arch string) []*pipeline.Pipe {
	var dmgURL, dmgFileName string

	latest, err := utils.GetLatestReleaseTag("influxdata", "telegraf")
	if err != nil {
		latest = "v1.33.1" // Default
	}
	latest = latest[1:]

	// Set the download URL and file name based on architecture
	if arch == "arm64" {
		dmgURL = "https://dl.influxdata.com/telegraf/releases/telegraf-" + latest + "_darwin_arm64.dmg"
		dmgFileName = "telegraf-" + latest + "_darwin_arm64.dmg"
	} else {
		dmgURL = "https://dl.influxdata.com/telegraf/releases/telegraf-" + latest + "_darwin_amd64.dmg"
		dmgFileName = "telegraf-" + latest + "_darwin_amd64.dmg"
	}

	volumeName := "/Volumes/Telegraf"

	pipes := []*pipeline.Pipe{
		{
			Name: "Downloading Telegraf DMG",
			Cmd:  exec.Command("curl", "-L", dmgURL, "-o", dmgFileName),
		},
		{
			Name: "Mounting DMG",
			Cmd:  exec.Command("hdiutil", "attach", dmgFileName),
		},
		{
			Name: "Moving telegraf app to /Applications",
			Cmd:  exec.Command("cp", "-R", volumeName+"/Telegraf.app", "/Applications/"),
		},
		{
			Name: "Copying telegraf binary to /usr/local/bin",
			Cmd:  exec.Command("cp", volumeName+"/Telegraf.app/Contents/Resources/usr/bin/telegraf", "/usr/local/bin/"),
		},
		{
			Name: "Detaching DMG",
			Cmd:  exec.Command("hdiutil", "detach", volumeName),
		},
	}

	return pipes

}

func DarwinUninstallPipes(sysinfo sysinfo.SysInfo) []*pipeline.Pipe {
	var pipes []*pipeline.Pipe
	pkgMngr := sysinfo.PkgMngr

	if pkgMngr == "brew" {
		pipes = BrewUninstallPipes()
	} else {
		pipes = macDmgUninstallPipes()
	}

	return pipes
}

func macDmgUninstallPipes() []*pipeline.Pipe {
	appPath := "/Applications/Telegraf.app"

	pipes := []*pipeline.Pipe{
		{
			Name: "Removing Telegraf From Applications",
			Cmd:  exec.Command("rm", "-rf", appPath),
		},
		{
			Name: "Removing Telegraf Binary",
			Cmd:  exec.Command("rm", "/usr/local/bin/telegraf"),
		},
	}

	return pipes
}
