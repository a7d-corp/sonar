/*
Copyright © 2021 Simon Weald

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package createconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glitchcrab/sonar/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dirPerms              os.FileMode = 0755
	filePerms             os.FileMode = 0600
	defaultConfigFileName             = "sonar.yaml"
	defaultConfigFilePath             = ".config/sonar"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Annotations: map[string]string{
			"skip-init-config": "true",
		},
		Use:   "create",
		Short: "Creates a Sonar config file if one does not already exist",
		RunE:  runCreateCommand,
	}

	return command
}

func runCreateCommand(cmd *cobra.Command, args []string) error {
	// Attempt to find an existing config file
	configFilePath, _ := config.FindConfigFile()

	if configFilePath != "" {
		log.Infof("config file already exists at %s", configFilePath)
		return nil
	}

	// Find the user's home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("could not establish user's home dir")
		return err
	}

	path := filepath.Join(homeDir, defaultConfigFilePath)

	// Create the config directory if it doesn't exist
	err = os.MkdirAll(path, dirPerms)
	if err != nil {
		return err
	}

	fullFilePath := filepath.Join(path, defaultConfigFileName)

	// Create the config file with default contents
	err = os.WriteFile(fullFilePath, []byte(defaultConfig), filePerms)
	if err != nil {
		return err
	}

	log.Infof("created config file at %s", fullFilePath)

	return nil
}
