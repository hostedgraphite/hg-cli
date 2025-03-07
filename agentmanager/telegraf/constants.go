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

var ServiceDetails = map[string]map[string]string{
	"windows": {
		"configPath": "C:\\Program Files\\InfluxData\\telegraf\\telegraf.conf",
		"serviceCmd": "C:\\Program Files\\InfluxData\\telegraf\\telegraf.exe",
		"startCmd":   "C:\\Program Files\\InfluxData\\telegraf\\telegraf.exe --service start",
		"restartCmd": "C:\\Program Files\\InfluxData\\telegraf\\telegraf.exe --service stop (then start)",
	},
	"linux": {
		"configPath": "/etc/telegraf/telegraf.conf",
		"confDir":    "/etc/telegraf",
		"serviceCmd": "/usr/bin/telegraf",
		"startCmd":   "sudo service telegraf start",
		"restartCmd": "sudo service telegraf restart",
	},
	"darwin": {
		"configPathAmd": "/usr/local/etc/telegraf.conf",
		"configDirAmd":  "/usr/local/etc",
		"configPathArm": "/opt/homebrew/etc/telegraf.conf",
		"confDirArm":    "/opt/homebrew/etc",
		"startCmd":      "brew services start telegraf",
		"serviceCmd":    "telegraf",
		"restartCmd":    "brew services restart telegraf",
	},
}
