package utils

import (
	"fmt"
	"net/http"
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
