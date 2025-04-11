package agentmanager

import (
	"strings"

	"github.com/hostedgraphite/hg-cli/agentmanager/otel"
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf"
	"github.com/hostedgraphite/hg-cli/sysinfo"
)

func GetAgent(agentName string) Agent {
	switch strings.ToLower(agentName) {
	case "telegraf":
		return &telegraf.Telegraf{}
	default:
		return nil
	}
}

func NewAgent(agentName string, options map[string]interface{}, sysInfo sysinfo.SysInfo) Agent {
	switch strings.ToLower(agentName) {
	case "telegraf":
		return telegraf.NewTelegrafAgent(options, sysInfo)
	case "otel", "opentelemetry":
		return otel.NewOtelAgent(options, sysInfo)
	default:
		return nil
	}
}
