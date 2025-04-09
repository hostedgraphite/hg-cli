package agents

import "github.com/charmbracelet/huh"

type AgentsFieldViews interface {
	InstallView(apikey, selectedInstall string, selectedPlugins []string) (*huh.Group, error)
	UninstallView(confirm bool) (*huh.Group, error)
	UpdateApiKeyView(apikey, path, defaultPath string) (*huh.Group, error)
}
