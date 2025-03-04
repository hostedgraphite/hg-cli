package telegraf

import (
	"time"

	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf/apiupdater"
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf/installers"
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf/uninstallers"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

type Telegraf struct {
	apikey  string
	sysinfo sysinfo.SysInfo
	options map[string]interface{}
	updates chan<- string
}

func (t *Telegraf) Install(apikey string, sysinfo sysinfo.SysInfo, options map[string]interface{}, updates chan<- string) error {
	t.apikey = apikey
	t.sysinfo = sysinfo
	t.options = options
	t.updates = updates

	var err error

	t.updates <- "Installing Telegraf Agent"
	err = installers.TelegrafAgentInstall(t.sysinfo, t.updates)
	if err != nil {
		updates <- "error installing Agent: " + err.Error()
	}
	time.Sleep(1 * time.Second)

	plugins := t.options["plugins"].([]string)

	time.Sleep(1 * time.Second)

	configPath := GetConfigPath(t.sysinfo.Os, t.sysinfo.Arch)
	telegrafCmd := ServiceDetails[t.sysinfo.Os]["serviceCmd"]

	updates <- "Installing Telegraf Plugin"
	err = installers.TelegrafPluginInstall(configPath, telegrafCmd, plugins, t.sysinfo)
	if err != nil {
		t.updates <- "error installing plugins: " + err.Error()
	}
	time.Sleep(1 * time.Second)

	updates <- "Updating Telegraf Graphite"
	err = installers.TelegrafGraphiteUpdate(t.apikey, configPath)
	if err != nil {
		t.updates <- "error installing plugins: " + err.Error()
	}

	t.updates <- "Completed Telegraf Agent Installation"

	return err
}

func (t *Telegraf) Uninstall(sysinfo sysinfo.SysInfo, updates chan<- string) error {
	var err error
	t.sysinfo = sysinfo
	t.updates = updates

	t.updates <- "Uninstalling Telegraf Agent"
	err = uninstallers.TelegrafUninstall(t.sysinfo, updates)
	if err != nil {
		t.updates <- "error uninstalling Agent: "
	}
	time.Sleep(1 * time.Second)

	t.updates <- "Deleting Additional Telegraf Files"
	err = uninstallers.TelegrafDeleteFiles(t.sysinfo, updates)
	if err != nil {
		t.updates <- "error deleting files: "
	}
	time.Sleep(1 * time.Second)

	t.updates <- "Completed Telegraf Agent Uninstall"
	return err
}

func (t *Telegraf) UpdateApiKey(apikey string, options map[string]interface{}, updates chan<- string) error {
	var err error
	t.apikey = apikey
	t.options = options
	t.updates = updates

	config := t.options["config"].(string)
	time.Sleep(1 * time.Second)

	updates <- "Updating Telegraf API Key"
	err = apiupdater.UpdateFile(t.apikey, config)
	if err != nil {
		updates <- "error updating API Key: "
	}
	time.Sleep(1 * time.Second)

	updates <- "Completed Telegraf API Key Update"
	return err
}
