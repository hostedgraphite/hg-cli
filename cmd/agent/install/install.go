package install

import (
	//  "fmt"
	"context"
	"fmt"

	"github.com/hostedgraphite/hg-cli/agentmanager"
	"github.com/hostedgraphite/hg-cli/agentmanager/telegraf"
	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
	"github.com/hostedgraphite/hg-cli/formatters"
	"github.com/hostedgraphite/hg-cli/sysinfo"

	// windows color support
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

func InstallCmd(sysinfo sysinfo.SysInfo) *cobra.Command {
	var (
		completed   bool
		installType string
		apikey      string
		agentName   string
		plugins     []string
	)

	cmd := &cobra.Command{
		Use:           "install <agent>",
		Short:         "Installing a monitoring agent.",
		Long:          "Install a moniting agent. Use --custom for custom installation",
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Check if the --list flag is added, which is a global flag
			list, _ := cmd.Flags().GetBool("list")
			if list {
				return nil
			}

			err := validateArgs(args, plugins, installType)
			if err != nil {
				return err
			}

			agentName = args[0]
			completed = true

			// Kind of wierd but, if the flags are marked required outisde of the PreRunE
			// the requiremet error will also populate when the --list flag is added.
			cmd.MarkFlagRequired("apikey")
			cmd.MarkFlagRequired("install")

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !completed {
				return nil
			}

			err := execute(apikey, installType, agentName, plugins, sysinfo)

			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&installType, "install", "", "The installation type (custom, default)")
	cmd.Flags().StringVar(&apikey, "apikey", "", "Your Hosted Graphite API key (required)")
	cmd.Flags().StringSliceVar(&plugins, "plugins", []string{}, "The plugins to install")

	return cmd
}

func validateArgs(args, plugins []string, installType string) error {

	if len(args) == 0 || !utils.ValidateAgent(args[0]) {
		return fmt.Errorf("no agent specified or agent not supported; see 'cli agent -l' for compatible agents")
	}

	if installType != "custom" && installType != "default" {
		return fmt.Errorf("the install type is not supported. Please select either custom or default")
	}

	if installType == "custom" {
		if len(plugins) == 0 {
			return fmt.Errorf("no plugins added, must include at least 1 plugin")
		}
	}

	return nil
}

func execute(apikey, installType, agentName string, plugins []string, sysinfo sysinfo.SysInfo) error {
	var err error
	var selectedPlugins []string

	agent := agentmanager.GetAgent(agentName)
	updates := make(chan string)

	if installType == "default" {
		selectedPlugins = telegraf.DefaultTelegrafPlugins
	} else {
		selectedPlugins = plugins
	}

	options := map[string]interface{}{
		"plugins": selectedPlugins,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		defer close(updates)
		err = agent.Install(apikey, sysinfo, options, updates)
		if err != nil {
			updates <- "error installing agent: " + err.Error()
			ctx.Err()
			return
		}
		ctx.Done()
	}(ctx)

	err = spinner.New().
		Title("Updating In Progress...").
		Context(ctx).
		ActionWithErr(func(ctx context.Context) error {
			for msg := range updates {
				fmt.Print("\r")
				fmt.Print("\033[K")
				fmt.Println(msg)
			}
			return nil
		}).
		Run()
	if err != nil {
		return err
	}

	summary := formatters.ActionSummary{
		Agent:    agentName,
		Success:  true,
		Action:   "Install",
		Plugins:  selectedPlugins,
		Config:   telegraf.GetConfigPath(sysinfo.Os, sysinfo.Arch),
		StartCmd: telegraf.ServiceDetails[sysinfo.Os]["startCmd"],
		Error:    "",
	}

	fmt.Println(formatters.GenerateCliSummary(summary))

	return err
}
