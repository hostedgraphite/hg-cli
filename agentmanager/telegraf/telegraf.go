package telegraf

import (
	"fmt"
	"time"

	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf/apiupdater"
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf/installers"
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf/uninstallers"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

type Telegraf struct {
	apikey          string
	sysinfo         sysinfo.SysInfo
	options         map[string]interface{}
	serviceSettings map[string]string
	updates         chan<- string
}

func NewTelegrafAgent(options map[string]interface{}, sysInfo sysinfo.SysInfo) *Telegraf {
	agent := &Telegraf{
		apikey:          options["apikey"].(string),
		sysinfo:         sysInfo,
		options:         options,
		serviceSettings: GetServiceSettings(sysInfo.Os, sysInfo.Arch, sysInfo.PkgMngr),
	}
	return agent
}

func (t *Telegraf) Install(apikey string, sysinfo sysinfo.SysInfo, options map[string]interface{}, updates chan<- string) error {
	t.apikey = apikey
	t.sysinfo = sysinfo
	t.options = options
	t.updates = updates

	var err error
	serviceSettings := GetServiceSettings(t.sysinfo.Os, t.sysinfo.Arch, t.sysinfo.PkgMngr)

	t.updates <- "Installing Telegraf Agent"
	err = installers.TelegrafAgentInstall(t.sysinfo, t.updates)
	if err != nil {
		updates <- "error installing Telegraf: " + err.Error()
		return err
	}
	time.Sleep(1 * time.Second)

	plugins := t.options["plugins"].([]string)

	time.Sleep(1 * time.Second)

	configPath := serviceSettings["configPath"]
	telegrafCmd := serviceSettings["serviceCmd"]

	updates <- "Configuring Telegraf Plugins"
	err = installers.TelegrafPluginInstall(configPath, telegrafCmd, plugins, t.sysinfo, t.updates)
	if err != nil {
		t.updates <- "error installing plugins: " + err.Error()
	}
	time.Sleep(1 * time.Second)

	updates <- "Updating Telegraf Graphite Output Config"
	err = installers.TelegrafGraphiteUpdate(t.apikey, configPath, t.sysinfo.Os, updates)
	if err != nil {
		t.updates <- "error updating telegraf config: " + err.Error()
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
		t.updates <- fmt.Sprintf("error uninstalling Agent: %v", err)
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

func (t *Telegraf) UpdateApiKey(apikey string, sysinfo sysinfo.SysInfo, options map[string]interface{}, updates chan<- string) error {
	var err error
	t.apikey = apikey
	t.options = options
	t.updates = updates
	t.sysinfo = sysinfo

	config := t.options["config"].(string)
	time.Sleep(1 * time.Second)

	updates <- "Updating Telegraf API Key"
	err = apiupdater.UpdateFile(t.apikey, config, t.sysinfo.Os)
	if err != nil {
		updates <- "error updating API Key: "
	}
	time.Sleep(1 * time.Second)

	updates <- "Completed Telegraf API Key Update"
	return err
}
