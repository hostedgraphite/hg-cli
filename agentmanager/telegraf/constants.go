package telegraf

var DefaultTelegrafPlugins = []string{
	"cpu",
	"disk",
	"diskio",
	"kernel",
	"mem",
	"processes",
	"swap",
	"system",
}

var ServiceDetails = map[string]map[string]map[string]string{
	"windows": {
		"default": {
			"configPath":  "C:\\Program Files\\InfluxData\\telegraf\\telegraf.conf",
			"serviceCmd":  "C:\\Program Files\\InfluxData\\telegraf\\telegraf.exe",
			"startHint":   `& "C:\Program Files\InfluxData\telegraf\telegraf.exe" --service-name telegraf service start`,
			"restartHint": `& "C:\Program Files\InfluxData\telegraf\telegraf.exe" --service-name telegraf service stop (then start)`,
		},
	},
	"linux": {
		"default": {
			"configPath":  "/etc/telegraf/telegraf.conf",
			"serviceCmd":  "telegraf",
			"startHint":   "sudo service telegraf start",
			"restartHint": "sudo service telegraf restart",
		},
		"brew": {
			"configPath":  "/home/linuxbrew/.linuxbrew/etc/telegraf.conf",
			"serviceCmd":  "telegraf",
			"startHint":   "brew service telegraf start",
			"restartHint": "brew service telegraf restart",
		},
	},
	"darwin": {
		"amd64": {
			"configPath":  "/usr/local/etc/telegraf.conf",
			"serviceCmd":  "telegraf",
			"startHint":   "brew services start telegraf",
			"restartHint": "brew services restart telegraf",
		},
		"arm64": {
			"configPath":  "/opt/homebrew/etc/telegraf.conf",
			"serviceCmd":  "telegraf",
			"startHint":   "brew services start telegraf",
			"restartHint": "brew services restart telegraf",
		},
	},
}
