package agent

import (
	"hg-cli/agentmanager/utils"
	"hg-cli/cmd/agent/apiupdater"
	"hg-cli/cmd/agent/install"
	"hg-cli/cmd/agent/uninstall"
	"hg-cli/sysinfo"

	"github.com/spf13/cobra"
)

func AgentCmd(sysinfo sysinfo.SysInfo) *cobra.Command {
	var listAgents bool

	cmd := &cobra.Command{
		Use:           "agent <command>",
		Short:         "Mange monitoring agents",
		Long:          "Install, Update, or uninstall monitoring agents.",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			if listAgents {
				utils.ShowAvailableAgents()
				return nil
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Since RunE will always run if there are no subcommands;
			// agent <install, update> ..., we need to exit in order
			// to prevent the help menu from loading.
			if listAgents {
				return nil
			}
			return cmd.Help()
		},
	}

	cmd.AddCommand(install.InstallCmd(sysinfo))
	cmd.AddCommand(uninstall.UninstallCmd(sysinfo))
	cmd.AddCommand(apiupdater.ApiUpdateCmd(sysinfo))
	cmd.PersistentFlags().BoolVarP(&listAgents, "list", "l", false, "List Available Agents")

	return cmd
}
