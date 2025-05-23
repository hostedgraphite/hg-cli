package uninstall

import (
	"fmt"

	"github.com/hostedgraphite/hg-cli/agentmanager"
	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/formatters"
	"github.com/hostedgraphite/hg-cli/pipeline"
	"github.com/hostedgraphite/hg-cli/sysinfo"
	cliUtils "github.com/hostedgraphite/hg-cli/utils"

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

			list, _ := cmd.Flags().GetBool("list")
			if list {
				return nil
			}

			err := validateArgs(args)
			if err != nil {
				return err
			}
			agentName = args[0]
			// Validate if the cmd requires sudo
			if cliUtils.AgentRequiresSudo(sysinfo.Os, "uninstall", sysinfo.PkgMngr, agentName) && !sysinfo.SudoPerm {
				return fmt.Errorf("this cmd requires admin privileges, please run as root")
			}
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

func execute(agentName string, sysInfo sysinfo.SysInfo) error {
	var err error
	var summary formatters.SummaryContent

	agent := agentmanager.NewAgent(agentName, nil, sysInfo)
	updates := make(chan *pipeline.Pipe)
	uninstallPipeline, err := agent.UninstallPipeline(updates)
	if err != nil {
		return err
	}

	runner := pipeline.NewRunner(
		uninstallPipeline,
		true,
		updates,
	)

	err = runner.Run()
	if err != nil {
		return err
	}

	data := formatters.ActionSummary{
		Agent:   agentName,
		Success: true,
		Action:  "Uninstall",
	}

	switch agentName {
	case "telegraf":
		summary = &formatters.TelegrafSummary{
			ActionSummary: data,
		}
	case "otel":
		summary = &formatters.OtelContribSummary{
			ActionSummary: data,
		}
	}

	fmt.Println(formatters.GenerateCliSummary(summary))

	return err
}
