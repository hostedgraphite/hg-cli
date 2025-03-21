package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
)

var agents = []string{"telegraf"}

func ShowAvailableAgents() {
	fmt.Println("Available agent: ")
	for _, agent := range agents {
		fmt.Println("- " + agent)
	}
}

func ValidateAgent(agent string) bool {
	return slices.Contains(agents, agent)
}

func UpdateConfigBlock(fullConfig, confBlock string, updates map[string]string) (string, error) {
	configRegex := regexp.MustCompile(confBlock)
	configBlock := configRegex.FindString(fullConfig)

	updatedBlock := configBlock

	if configBlock == "" {
		return "", fmt.Errorf("error: no matching graphite configuration found")
	}

	for regexPattern, replacement := range updates {
		re := regexp.MustCompile(regexPattern)
		updatedBlock = re.ReplaceAllString(updatedBlock, replacement)
	}

	updatedConfig := configRegex.ReplaceAllString(fullConfig, updatedBlock)

	return updatedConfig, nil
}

func GetLatestReleaseTag(repo_org string, repo_name string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repo_org, repo_name)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("unable to get latest release tag: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	return release.TagName, nil
}
