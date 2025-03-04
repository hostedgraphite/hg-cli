package utils

import (
	"fmt"
	"slices"
)

var agents = []string{"telegraf", "collectd", "prom"}

func ShowAvailableAgents() {
	fmt.Println("Available agent: ")
	for _, agent := range agents {
		fmt.Println("- " + agent)
	}
}
func ValidateAgent(agent string) bool {
	return slices.Contains(agents, agent)
}
