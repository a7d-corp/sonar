/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	image             string
	networkPolicy     bool
	podSecurityPolicy bool
	privileged        bool

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "create will deploy a debug container to the current context",
		Long: `Create will create a debugging deployment and all supporting resources to
the provided kubectl context (or the current context if none is provided).
Sonar assumes that your context has the required privileges.

All flags are optional as sane defaults are provided.

Image names may be provided with or without a tag; if no tag is detected
then the 'latest' tag is automatically used.`,
		PreRun: validateFlags,
		Run:    createSonarDeployment,
	}
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&image, "image", "i", "busybox", "image name (e.g. glitchcrab/ubuntu-debug:latest)")
	createCmd.Flags().BoolVar(&networkPolicy, "network-policy", false, "create NetworkPolicy (default \"false\")")
	createCmd.Flags().BoolVar(&podSecurityPolicy, "podsecuritypolicy", false, "create PodSecurityPolicy (default \"false\")")
	createCmd.Flags().BoolVar(&privileged, "privileged", false, "run the container as root (assumes userID of 0) (default \"false\")")
}

func validateFlags(cmd *cobra.Command, args []string) {
	if image == "" {
		log.Fatal("Image name for debugging container must be provided")
	}
}

func createSonarDeployment(cmd *cobra.Command, args []string) {
	fmt.Printf("name: %s\n", name)
	fmt.Printf("namespace: %s\n", namespace)
	fmt.Printf("image: %s\n", image)
}
