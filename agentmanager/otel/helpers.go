package otel

import (
	"fmt"
	"os"
)

func GetServiceSettings(os, arch, pkgmngr string) map[string]string {
	var settings map[string]string

	switch os {
	case "windows":
		settings = ServiceDetails[os]["default"]
	case "linux":
		if pkgmngr == "brew" {
			settings = ServiceDetails[os][pkgmngr]
		} else {
			settings = ServiceDetails[os]["default"]
		}
	case "darwin":
		switch arch {
		case "amd64":
			settings = ServiceDetails[os][arch]
		default:
			settings = ServiceDetails[os][arch]
		}
	}

	return settings
}

func ValidateFilePath(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("error verifying path: %v", err)
	}

	// check that it's a file
	if info.IsDir() {
		return fmt.Errorf("error verifying path: %v is a directory", filePath)
	}

	return nil
}
