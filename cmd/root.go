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
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	nameStub = "sonar"
)

var (
	cfgFile   string
	name      string
	namespace string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "sonar",
		Short: "Sonar is a tool for deploying a standalone debugging container to a Kubernetes cluster.",
		Long: `Sonar is used to create a Kubernetes deployment with a debug
container for troubleshooting cluster issues.

The deployment can be customised to a certain extent in
order to suit the target cluster configuration.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&name, "name", "debug", "deployment name")
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "namespace")
}

func initConfig() {
	// If the user hasn't provided a deployment name then inform them that
	// we are using the default. Else we validate that the name looks sane.
	if !rootCmd.PersistentFlags().Lookup("name").Changed {
		log.Infof("No name provided; defaulting name to: %s-%s", nameStub, name)
	} else {
		// validate provided name
		// max 253 chars, only alphanumeric, -. only, start/end alphanumeric
	}

	// If the user hasn't provided a deployment namespace then inform them that
	// we are using the default. Else we validate that the namespace looks sane.
	if !rootCmd.PersistentFlags().Lookup("namespace").Changed {
		log.Infof("No namespace provided, deploying to: %s", namespace)
	} else {
		// validate provided namespace
	}
}
