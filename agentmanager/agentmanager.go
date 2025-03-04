package agentmanager

import (
	"strings"

	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf"
)

func GetAgent(agentName string) Agent {
	switch strings.ToLower(agentName) {
	case "telegraf":
		return &telegraf.Telegraf{}
	default:
		return nil
	}
}
