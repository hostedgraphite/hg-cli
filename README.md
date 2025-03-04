# HG-CLI
Current Flags:
* --apikey (required)
* --action (install, update-apikey)
* --agent (telegraf)
* --installType (default, custom)
* --plugins (custom install)
* --config (config path for updating apikey)

To run in TUI mode:
go run . 
or use the executable

To run CLI mode - Install:
```<go run . or .exe> --apikey <apikey> --action install --agent Telegraf```

That will install the default plugins.

To run cusom install add the flags:

```--installType custom```

```--plugins cpu,disk,...```

To make an apikey update:

```<go run . or .exe> --apikey <apikey> --action update-apikey --agent Telegraf --config <config path>```


### Currently Set to Sandbox Env and API validation is active, so you'll need to use a sandbox apikey.# hg-cli
