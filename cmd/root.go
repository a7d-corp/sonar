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
	"github.com/spf13/cobra"
)

var (
	kubeConfig  string
	kubeContext string
	name        string
	namespace   string
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
		PersistentPreRunE: initRootConfig,
	}

	// Add global flags
	root.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "", "absolute path to kubeconfig file (default: '/home/$user/.kube/config')")
	root.PersistentFlags().StringVar(&kubeContext, "context", "", "context to use")
	root.PersistentFlags().StringVarP(&name, "name", "N", "", "resource name (max 50 characters) (automatically prepended with 'sonar-')")
	root.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace to operate in")

	root.AddCommand(version.NewCommand())

	// Add subcommands
	root.AddCommand(
		create.NewCommand(),
		destroy.NewCommand(),
		exec.NewCommand(),
		ls.NewCommand(),
		version.NewCommand(),
	)

	return root
}

func initRootConfig(root *cobra.Command, args []string) error {
	// Skip config initialisation for commands which do not need it.
	if root.Annotations["skip-init-config"] == "true" {
		return nil
	}

	// Create config struct.
	globals := config.Globals{
		KubeConfig:  kubeConfig,
		KubeContext: kubeContext,
		Labels:      make(map[string]string),
		Name:        name,
		Namespace:   namespace,
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
