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
	"fmt"
	"regexp"

	"github.com/glitchcrab/sonar/service/k8sclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	nameMaxLength  = 50
	nameRegex      = "^[a-zA-Z0-9-.]*$"
	nameStub       = "sonar"
	namespaceRegex = "^[a-zA-Z0-9-]*$"
)

var (
	fullName    string
	kubeConfig  string
	kubeContext string
	labels      = map[string]string{
		"created-by": "sonar",
		"owner":      "sonar",
	}
	name      string
	namespace string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
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
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "", "absolute path to kubeconfig file (default: '/home/$user/.kube/config')")
	rootCmd.PersistentFlags().StringVar(&kubeContext, "context", "", "context to use")
	rootCmd.PersistentFlags().StringVarP(&name, "name", "N", "", "resource name (max 50 characters) (automatically prepended with 'sonar-')")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace to operate in")
}

func initConfig() {
	// If the user has provided a name then validate that it looks sane.
	if rootCmd.PersistentFlags().Lookup("name").Changed {
		// Restrict deployment name to 50 characters. 50 is a relatively
		// arbitrary choice of length, but it should be sufficient.
		if len(name) > nameMaxLength {
			log.Fatal("deployment name must be 50 characters or less")
		}

		// Validate the provided name is suitable for a Kubernetes resource name.
		ok, _ := regexp.MatchString(nameRegex, name)
		if !ok {
			log.Fatal("deployment name can only contain alphanumeric characters, hyphens and periods")
		}
	}

	// Prepend provided name with 'sonar-' for ease of identifying Sonar deployments.
	if name != "" {
		fullName = fmt.Sprintf("%s-%s", nameStub, name)
	} else {
		fullName = nameStub
	}

	// Add the provided name to the labels map for tagging generated resources
	labels["name"] = fullName

	// If the user has provided a namespace then validate that it looks sane.
	if rootCmd.PersistentFlags().Lookup("namespace").Changed {
		// Validate the provided namespace is suitable for a Kubernetes namespace.
		ok, _ := regexp.MatchString(namespaceRegex, namespace)
		if !ok {
			log.Fatal("namespaces can only contain alphanumeric characters and hyphens")
		}
	}

	// If the namespace was not provided, get the namespace from the context
	var err error
	if namespace == "" {
		namespace, err = k8sclient.GetNamespace(kubeConfig, kubeContext)
		if err != nil {
			log.Fatal(err) // TODO: better logging
		}
	}
}
