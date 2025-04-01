package otel

import (
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

type Otel struct {
	apikey          string
	sysinfo         sysinfo.SysInfo
	options         map[string]interface{}
	serviceSettings map[string]string
	updates         chan<- string
}

func NewOtelAgent(options map[string]interface{}, sysInfo sysinfo.SysInfo) *Otel {
	apikey, ok := options["apikey"].(string)
	if !ok {
		apikey = ""
	}

	agent := &Otel{
		apikey:          apikey,
		sysinfo:         sysInfo,
		options:         options,
		serviceSettings: GetServiceSettings(sysInfo.Os, sysInfo.Arch, sysInfo.PkgMngr),
	}

	return agent
}
