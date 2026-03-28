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
package version

import (
	"github.com/glitchcrab/sonar/pkg/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "Prints the version of Sonar",
		Annotations: map[string]string{
			"skip-init-config": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			version := project.Version()
			log.Infof("version: %s\n", version)
			return nil
		},
	}

	return command
}
