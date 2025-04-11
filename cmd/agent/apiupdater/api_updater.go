package apiupdater

import (
	"fmt"

	"github.com/hostedgraphite/hg-cli/agentmanager"
	"github.com/hostedgraphite/hg-cli/agentmanager/otel"
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf"
	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/formatters"
	"github.com/hostedgraphite/hg-cli/pipeline"
	"github.com/hostedgraphite/hg-cli/sysinfo"
	cliUtils "github.com/hostedgraphite/hg-cli/utils"

	"github.com/spf13/cobra"
)

func ApiUpdateCmd(sysinfo sysinfo.SysInfo) *cobra.Command {
	var agentName, apikey, path string
	var completed bool

	cmd := &cobra.Command{
		Use:   "update-apikey <agent>",
		Short: "Update the API key for a monitoring agent.",
		Long:  "Update the API key for a monitoring agent.",
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

			// Otel will require sudo permissions to update the config file
			if agentName == "otel" && cliUtils.ActionRequiresSudo(sysinfo.Os, "update", sysinfo.PkgMngr) && !sysinfo.SudoPerm {
				return fmt.Errorf("this cmd requires admin privileges, please run as root")
			}

			completed = true

			cmd.MarkFlagRequired("apikey")
			cmd.MarkFlagRequired("config")

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !completed {
				return nil
			}

			err := execute(apikey, agentName, path, sysinfo)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&apikey, "api-key", "", "Your Hosted Graphite API key (required)")
	cmd.Flags().StringVar(&path, "config", "", "The path to the agent configuration file")

	return cmd
}

func validateArgs(args []string) error {
	if len(args) == 0 || !utils.ValidateAgent(args[0]) {
		return fmt.Errorf("no agent specified or agent not supported; see 'cli agent -l' for compatible agents")
	}
	return nil
}

func execute(apikey, agentName, path string, sysInfo sysinfo.SysInfo) error {
	var err error
	var serviceSettings map[string]string
	var summary formatters.SummaryContent

	options := map[string]interface{}{
		"config": path,
		"apikey": apikey,
	}
	agent := agentmanager.NewAgent(agentName, options, sysInfo)
	updates := make(chan *pipeline.Pipe)
	switch agentName {
	case "telegraf":
		serviceSettings = telegraf.GetServiceSettings(sysInfo.Os, sysInfo.Arch, sysInfo.PkgMngr)
		summary = &formatters.TelegrafSummary{
			ActionSummary: formatters.ActionSummary{
				Agent:      agentName,
				Success:    true,
				Action:     "Update Api Key",
				Config:     path,
				RestartCmd: serviceSettings["restartHint"],
			},
		}
	case "otel":
		serviceSettings = otel.GetServiceSettings(sysInfo.Os, sysInfo.Arch, sysInfo.PkgMngr)
		summary = &formatters.OtelContribSummary{
			ActionSummary: formatters.ActionSummary{
				Agent:      agentName,
				Success:    true,
				Action:     "Update Api Key",
				Config:     path,
				RestartCmd: serviceSettings["restartHint"],
			},
		}
	}
	updateApikeyPipeline, err := agent.UpdateApiKeyPipeline(updates)
	if err != nil {
		return err
	}

	runner := pipeline.NewRunner(
		updateApikeyPipeline,
		true,
		updates,
	)

	err = runner.Run()
	if err != nil {
		return err
	}

	fmt.Println(formatters.GenerateCliSummary(summary))

	return err
}
