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
	"context"

	"github.com/glitchcrab/sonar/cmd/create"
	"github.com/glitchcrab/sonar/cmd/destroy"
	"github.com/glitchcrab/sonar/cmd/exec"
	"github.com/glitchcrab/sonar/cmd/ls"
	"github.com/glitchcrab/sonar/cmd/version"
	"github.com/glitchcrab/sonar/internal/app"
	"github.com/glitchcrab/sonar/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultConfigFilename = "sonar"
	configFile            string
	kubeConfig            string
	kubeContext           string
	name                  string
	namespace             string
	v                     *viper.Viper
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "sonar",
		Short: "Sonar deploys a debugging container to a Kubernetes cluster.",
		Long: `Sonar is used to create a Kubernetes deployment with a debug container
for troubleshooting cluster issues.

The deployment can be customised to a certain extent in order to
suit the target cluster configuration.

Global flags:

--kubeconfig (default: '/home/$user/.kube/config')

Absolute path to the kubeconfig file to use.
If left blank, Sonar will read the KUBECONFIG environment variable.
If that is not set, Sonar will use the default kubeconfig file location.

--context (default: current context in kube config)

Name of the kubernetes context to use.

--name

Name given to all the created resources. If provided then this will
be automatically prepended with 'sonar-' for ease of identification.
For example, a provided name of 'test' will result in a deployment named
'sonar-test'. Provided name can be a max of 50 characters. If no name is
provided then resource names will start with 'sonar-'.

--namespace (default: 'default')

Namespace to deploy resources to.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate user-prvided config.
			var err error
			err = initRootConfig(cmd, args, v)
			if err != nil {
				return err
			} else {
				return nil
			}
		},
	}

	// Add global flags
	root.PersistentFlags().StringVar(&configFile, "config", "", "path to Sonar config file")
	root.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "", "absolute path to kubeconfig file (default: '/home/$user/.kube/config')")
	root.PersistentFlags().StringVar(&kubeContext, "context", "", "context to use")
	root.PersistentFlags().StringVarP(&name, "name", "N", "", "resource name (max 50 characters) (automatically prepended with 'sonar-')")
	root.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace to operate in")

	root.AddCommand(version.NewCommand())

	// Wire in Viper.
	v, _ = initViperConfig(root)

	// Add subcommands
	root.AddCommand(
		create.NewCommand(v),
		destroy.NewCommand(),
		exec.NewCommand(),
		ls.NewCommand(),
		version.NewCommand(),
	)

	return root
}

func initRootConfig(root *cobra.Command, args []string, v *viper.Viper) error {
	// Skip config initialisation for commands which do not need it.
	if root.Annotations["skip-init-config"] == "true" {
		return nil
	}

	// Create config struct.
	globals := config.Globals{
		KubeConfig:  kubeConfig,
		KubeContext: kubeContext,
		Labels:      make(map[string]string),
		Name:        v.GetString("name"),
		Namespace:   v.GetString("namespace"),
	}

	// Validate user-provided config values.
	if err := config.ValidateGlobalConfig(&globals); err != nil {
		return err
	}

	// Instantiate an App struct.
	app := &app.App{Globals: globals}

	appKey := app.RetrieveAppKey()

	// Add the App struct to the command's context.
	ctx := context.WithValue(root.Context(), appKey, app)
	root.SetContext(ctx)

	return nil
}

// initViperConfig initialises a Viper instance and binds some Cobra flags.
func initViperConfig(root *cobra.Command) (*viper.Viper, error) {
	v := viper.New()

	// Use the provided config file.
	if configFile != "" {
		// Set the user-provided config file.
		v.SetConfigFile(configFile)
	} else {
		// Search for the config file in the user's home directory.
		v.SetConfigName(defaultConfigFilename)
		v.AddConfigPath("$HOME/.config")
		v.AddConfigPath("$HOME/.config/sonar")
	}

	// Attempt to read a config file.
	if err := v.ReadInConfig(); err != nil {
		// Ignore file not found errors, but bail on any other error (such as parsing failures).
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		log.Info("no config file found")
	} else {
		log.Infof("using config file: %s", v.ConfigFileUsed())
	}

	// Bind some flags to Viper.
	if err := v.BindPFlag("name", root.PersistentFlags().Lookup("name")); err != nil {
		return nil, err
	}
	if err := v.BindPFlag("namespace", root.PersistentFlags().Lookup("namespace")); err != nil {
		return nil, err
	}

	return v, nil
}
