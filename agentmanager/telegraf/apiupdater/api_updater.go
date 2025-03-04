package apiupdater

import (
	"fmt"
	"os"
	"regexp"
)

func UpdateFile(apikey, filePath string) error {
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
