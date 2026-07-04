package config

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	// Potential file extensions
	extensions = [3]string{"yaml", "json", "toml"}
	// Default config file name
	configFileName = "sonar"
	// Default config file path
	configFilePaths = [2]string{".config", ".config/sonar"}
)

// FindConfigFile attempts to locate a config file at the expected locations
func FindConfigFile() (string, error) {
	// Find the user's home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Info("could not establish user's home dir")
		return "", err
	}

	// Check if the config file exists at the default paths with all possible extensions
	for _, ext := range extensions {
		for _, path := range configFilePaths {
			// Build the full path to the config file
			fullPath := fmt.Sprintf("%s.%s", filepath.Join(homeDir, path, configFileName), ext)

			// Check if the file exists
			_, err := os.Stat(fullPath)
			// Return if the file exists
			if err == nil {
				return fullPath, nil
			}
		}
	}

	// Return an error if the config file wasn't found
	return "", fmt.Errorf("no config file found")
}
