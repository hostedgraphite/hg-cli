package telegraf

import (
	"fmt"
	"os"
)

func GetConfigPath(os, arch string) string {
	var path string

	switch os {
	case "windows":
		path = ServiceDetails[os]["configPath"]
	case "linux":
		path = ServiceDetails[os]["configPath"]
	case "darwin":
		switch arch {
		case "amd64":
			path = ServiceDetails[os]["configPathAmd"]
		default:
			path = ServiceDetails[os]["configPathArm"]
		}
	}

	return path
}

func ValidateFilePath(filePath, osys, arch string, tui bool) error {

	if tui && filePath == "" {
		filePath = GetConfigPath(osys, arch)
	}

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
