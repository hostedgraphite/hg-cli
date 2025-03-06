package apiupdater

import (
	"fmt"
	"strings"

	"github.com/hostedgraphite/hg-cli/agentmanager"
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf"
	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/formatters"
	"github.com/hostedgraphite/hg-cli/sysinfo"

	"github.com/charmbracelet/huh/spinner"
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
			completed = true

			cmd.MarkFlagRequired("apikey")
			cmd.MarkFlagRequired("path")

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

	cmd.Flags().StringVar(&apikey, "apikey", "", "Your Hosted Graphite API key (required)")
	cmd.Flags().StringVar(&path, "config", "", "The path to the agent configuration file")

	return cmd
}

func validateArgs(args []string) error {
	if len(args) == 0 || !utils.ValidateAgent(args[0]) {
		return fmt.Errorf("no agent specified or agent not supported; see 'cli agent -l' for compatible agents")
	}
	return nil
}

func execute(apikey, agentName, path string, sysinfo sysinfo.SysInfo) error {
	var err error

	agent := agentmanager.GetAgent(agentName)
	updates := make(chan string)
	options := map[string]interface{}{
		"config": path,
	}

	go func() {
		defer close(updates)
		err = agent.UpdateApiKey(apikey, sysinfo, options, updates)
		if err != nil {
			updates <- "error updating api key"
			return
		}
	}()

	err = spinner.New().Title("Updating In Progress...").Action(func() {
		for msg := range updates {
			fmt.Println(msg)
			if strings.HasPrefix(msg, "error") || strings.HasPrefix(msg, "Completed") {
				break
			}
		}
	}).Run()

	if err != nil {
		return fmt.Errorf("error updating api key: %v", err)
	}

	summary := formatters.ActionSummary{
		Agent:    agentName,
		Success:  true,
		Action:   "Update Api Key",
		Config:   path,
		StartCmd: telegraf.ServiceDetails[sysinfo.Os]["restartCmd"],
	}

	fmt.Print(formatters.GenerateCliSummary(summary))

	return err
}
