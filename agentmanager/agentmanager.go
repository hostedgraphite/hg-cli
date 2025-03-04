package agentmanager

import (
	"hg-cli/agentmanager/telegraf"
	"strings"
)

func GetAgent(agentName string) Agent {
	switch strings.ToLower(agentName) {
	case "telegraf":
		return &telegraf.Telegraf{}
	default:
		return nil
	}
}
