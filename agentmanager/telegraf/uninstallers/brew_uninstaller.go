package uninstallers

import (
	"fmt"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
)

func BrewUninstall(updates chan<- string) error {
	var err error

	if err = utils.RunCommand("brew", []string{"services", "stop", "telegraf"}, updates); err != nil {
		return fmt.Errorf("error stopping telegraf service: %v", err)
	}

	if err = utils.RunCommand("brew", []string{"uninstall", "telegraf"}, updates); err != nil {
		return fmt.Errorf("error uninstalling telegraf service: %v", err)
	}
	return err
}
