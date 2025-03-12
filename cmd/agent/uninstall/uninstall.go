package uninstall

import (
	"fmt"
	"strings"

	"github.com/hostedgraphite/hg-cli/agentmanager"
	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/sysinfo"
	cliUtils "github.com/hostedgraphite/hg-cli/utils"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

func UninstallCmd(sysinfo sysinfo.SysInfo) *cobra.Command {
	var agentName string
	var completed bool

	cmd := &cobra.Command{
		Use:   "uninstall <agent>",
		Short: "Uninstall a monitoring agent.",
		Long:  "Uninstall a monitoring agent.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate if the cmd requires sudo
			if cliUtils.ActionRequiresSudo(sysinfo.Os, "uninstall") && !sysinfo.SudoPerm {
				return fmt.Errorf("this cmd requires admin privileges, please run as root")
			}

			list, _ := cmd.Flags().GetBool("list")
			if list {
				return nil
			}

			err := validateArgs(args)
			if err != nil {
				return err
			}
			agentName = args[0]
			completed = true
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !completed {
				return nil
			}

			err := execute(agentName, sysinfo)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func validateArgs(args []string) error {
	var err error

	if len(args) == 0 || !utils.ValidateAgent(args[0]) {
		return fmt.Errorf("no agent specified or agent not supported; see 'cli agent -l' for compatible agents")
	}

	return err
}

func execute(agentName string, sysinfo sysinfo.SysInfo) error {
	var err error

	agent := agentmanager.GetAgent(agentName)
	updates := make(chan string)

	go func() {
		defer close(updates)
		err = agent.Uninstall(sysinfo, updates)
		if err != nil {
			updates <- fmt.Sprintf("error uninstalling agent: %v", err)
			return
		}
	}()

	err = spinner.New().Title("Uninstall In Progress...").Action(func() {
		for msg := range updates {
			fmt.Print("\r")
			fmt.Print("\033[K")
			fmt.Println(msg)
			if strings.HasPrefix(msg, "error") || strings.HasPrefix(msg, "Completed") {
				break
			}
		}
	}).Run()

	if err != nil {
		return fmt.Errorf("error uninstalling agent: %v", err)
	}

	return err
}
