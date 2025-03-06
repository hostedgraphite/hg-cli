package agentmanager

import "github.com/hostedgraphite/hg-cli/sysinfo"

type Agent interface {
	Install(apikey string, sysinfo sysinfo.SysInfo, options map[string]interface{}, updates chan<- string) error
	Uninstall(sysinfo sysinfo.SysInfo, updates chan<- string) error
	UpdateApiKey(apikey string, sysinfo sysinfo.SysInfo, options map[string]interface{}, updates chan<- string) error
}
