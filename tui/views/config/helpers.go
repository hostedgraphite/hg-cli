package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed plugins.yaml
var pluginsYaml []byte

type PluginYaml struct {
	Plugins []string `yaml:"plugins"`
}

func LoadPlugins() (*PluginYaml, error) {
	var options PluginYaml
	err := yaml.Unmarshal(pluginsYaml, &options)
	if err != nil {
		return nil, fmt.Errorf("error loading yaml: %w", err)
	}

	return &options, nil
}
