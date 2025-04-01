package otel

var DefaultConfig = []string{
	"Receiver: hostmetrics",
	"Exporters: Carbon",
}

var ServiceDetails = map[string]map[string]map[string]string{
	"linux": {
		"default": {
			"configPath":  "/etc/otelcol-contrib/config.yaml",
			"startHint":   "sudo systemctl start otelcol-contrib",
			"restartHint": "sudo systemctl restart otelcol-contrib",
		},
	},
}
