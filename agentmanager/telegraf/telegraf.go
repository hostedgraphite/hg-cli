package telegraf

import (
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
	var apikey string
	// Apikey should have been verified by this point,
	// both in the TUI and the CLI. Only the installation and
	// udpate key should have the options with the apikey
	// so it can be set to empty string here.
	apikey, ok := options["apikey"].(string)
	if !ok {
		apikey = ""
	}

	agent := &Telegraf{
		apikey:          apikey,
		sysinfo:         sysInfo,
		options:         options,
		serviceSettings: GetServiceSettings(sysInfo.Os, sysInfo.Arch, sysInfo.PkgMngr),
	}
	return agent
}
