package cmd

import (
	"fmt"
	"hg-cli/cmd/agent"
	"hg-cli/styles"
	"hg-cli/sysinfo"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "hg-cli",
	Short:         "CLI to interact with Hosted Graphite",
	Long:          "CLI to interact with Hosted Graphite",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Add top level commands here
func init() {
	sysinfo, err := sysinfo.GetSystemInformation()
	if err != nil {
		fmt.Printf("error getting system information: %v", err)
		os.Exit(1)
	}
	rootCmd.AddCommand(TuiEnableCmd(sysinfo))
	rootCmd.AddCommand(agent.AgentCmd(sysinfo))
	rootCmd.SetUsageFunc(styles.CustomUsageFunc)
}

func Execute() {
	s := styles.DefaultStyles()
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		content := err.Error()
		fmt.Println(s.Error.Render(content))
	}
}
