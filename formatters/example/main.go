package main

import (
	"fmt"

	"github.com/hostedgraphite/hg-cli/formatters"
)

func main() {
	var summary1, summary2 formatters.SummaryContent
	options := make(map[string]interface{})
	options["plugins"] = []string{"cpu", "disk", "mem"}
	summary1 = &formatters.TelegrafSummary{
		ActionSummary: formatters.ActionSummary{
			Agent:    "Telegraf",
			Success:  true,
			Action:   "Install",
			Config:   "/etc/telegraf/telegraf.conf",
			StartCmd: "sudo service telegraf start",
		},
		Plugins: options["plugins"].([]string),
	}

	summary2 = &formatters.OtelContribSummary{
		ActionSummary: formatters.ActionSummary{
			Agent:    "OpenTelemetry",
			Success:  true,
			Action:   "Install",
			Config:   "/etc/otelcontrib/otelcontrib.conf",
			StartCmd: "sudo service otelcontrib start",
		},
		Receiver: "hostmetrics",
		Exporter: "carbon",
	}

	fmt.Println("formatters:", formatters.GenerateSummary(summary1, 80, 20))

	fmt.Println("formatters:", formatters.GenerateSummary(summary2, 80, 20))
}
