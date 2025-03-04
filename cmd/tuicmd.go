package cmd

import (
	// "fmt"
	"fmt"
	"hg-cli/sysinfo"
	"hg-cli/tui"
	"os"

	"github.com/spf13/cobra"
)

func TuiEnableCmd(sysinfo sysinfo.SysInfo) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Enable tui",
		Long:  "Enable tui",
		Run: func(cmd *cobra.Command, args []string) {
			if err := tui.StartTui(sysinfo); err != nil {
				fmt.Println("Error launching TUI:", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}
