package utils

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func ValidateAPIKey(apikey string) error {
	// We just need to know that the api key is valid so a quick query to any metric will work,
	// even if the metric doesn't exist.
	url := fmt.Sprintf("https://%s@api.hostedgraphite.com/api/v1/metric/search?pattern=test", apikey)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid API key, received status code: %d", resp.StatusCode)
	}

	return nil
}

func AgentRequiresSudo(os, action, pkgmngr, agent string) bool {
	needSudo := true

	if agent == "otel" {
		return true
	} else if agent == "telegraf" {
		if pkgmngr == "brew" {
			return false
		}
	}

	return needSudo
}

// Previous sudo checker
func ActionRequiresSudo(os, action, pkgmngr string) bool {
	if pkgmngr == "brew" {
		return false
	}
	action = strings.ToLower(action)
	needSudo := true
	sudoActions := []string{"install", "uninstall", "update"}

	switch os {
	case "darwin":
		// We'll install with brew
		needSudo = false
	case "linux":
		needSudo = slices.Contains(sudoActions, action)
	case "windows":
		needSudo = slices.Contains(sudoActions, action)
	default:
		needSudo = false
	}

	return needSudo
}
