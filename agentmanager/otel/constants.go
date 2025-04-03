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
	"windows": {
		"default": {
			"exePath":     "C:\\Program Files\\OpenTelemetry Collector Contrib\\otelcol-contrib.exe",
			"configPath":  "C:\\Program Files\\OpenTelemetry Collector Contrib\\config.yaml",
			"startHint":   "sc.exe start otelcol-contrib",
			"restartHint": "sc.exe restart otelcol-contrib",
		},
	},
}
