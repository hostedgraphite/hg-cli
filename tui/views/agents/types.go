package agents

import "github.com/charmbracelet/huh"

type AgentsFieldViews interface {
	InstallView() (*huh.Group, error)
	UninstallView() (*huh.Group, error)
	UpdateApiKeyView(defaultPath string) (*huh.Group, error)
}
