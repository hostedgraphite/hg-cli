package installers

import (
	"fmt"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
)

func BrewInstall(updates chan<- string) error {
	var err error
	if err = utils.RunCommand("brew", []string{"install", "telegraf"}, updates); err != nil {
		return fmt.Errorf("error installing telegraf service: %v", err)
	}
	return err
}
