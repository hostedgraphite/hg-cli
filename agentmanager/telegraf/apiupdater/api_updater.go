package apiupdater

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hostedgraphite/hg-cli/agentmanager/utils"
)

func UpdateFile(apikey, filePath, opersystem string) error {

	if opersystem == "linux" {
		err := linuxApiUpdater(apikey, filePath)
		if err != nil {
			return fmt.Errorf("error updating file: %v", err)
		}
		return nil
	}

	newPrefix := apikey + ".telegraf"
	re := regexp.MustCompile(`(?m)^\s*prefix\s*=\s*".*"`)
	currentConfig, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	updatedConfig := re.ReplaceAllString(string(currentConfig), fmt.Sprintf(`prefix = "%s"`, newPrefix))

	err = os.WriteFile(filePath, []byte(updatedConfig), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func linuxApiUpdater(apikey, filePath string) error {
	fullConfig, err := utils.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	graphiteBlock := `\[\[outputs\.graphite\]\](?:.|\s)*?\[\[`
	updates := map[string]string{
		`prefix\s*=\s*".*?"`: fmt.Sprintf(`prefix = "%s.telegraf"`, apikey),
	}
	updatedConfig, err := utils.UpdateConfigBlock(fullConfig, graphiteBlock, updates)
	if err != nil {
		return fmt.Errorf("error writing file:%v", err)
	}

	err = utils.WriteFile(filePath, updatedConfig)
	if err != nil {
		return fmt.Errorf("error writing file:%v", err)
	}

	return nil
}
