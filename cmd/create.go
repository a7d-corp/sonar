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
	imageRegex = "^[a-z0-9/.-]*[:][a-z0-9.-]*$"
)

var (
	image             string
	networkPolicy     bool
	podSecurityPolicy bool
	privileged        bool

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "create applies a debug deployment to a Kubernetes cluster",
		Long: `Create will attempt to create a debugging deployment and all supporting
resources in the provided kubectl context (or the current context if
none is provided). Sonar assumes that your context has the required
privileges to create the necessary resources.

All flags are optional as sane defaults are provided.

Image names may be provided with or without a tag; if no tag is detected
then the 'latest' tag is automatically used.`,
		Run: createSonarDeployment,
	}
)

func init() {
	rootCmd.AddCommand(createCmd)
	cobra.OnInitialize(validateFlags)

	createCmd.Flags().StringVarP(&image, "image", "i", "busybox:latest", "image name (e.g. glitchcrab/ubuntu-debug:latest)")
	createCmd.Flags().BoolVar(&networkPolicy, "network-policy", false, "create NetworkPolicy (default \"false\")")
	createCmd.Flags().BoolVar(&podSecurityPolicy, "podsecuritypolicy", false, "create PodSecurityPolicy (default \"false\")")
	createCmd.Flags().BoolVar(&privileged, "privileged", false, "run the container as root (assumes userID of 0) (default \"false\")")
}

func validateFlags() {
	// If the user hasn't provided an image name then inform them that
	// we are using the default. Else we validate the image tag.
	if !createCmd.Flags().Lookup("image").Changed {
		log.Infof("No image name provided, using: %s", image)
	} else {
		// Validate image to see if a tag has been provided; if not then
		// use latest. Does not validate full image name, just whether a
		// tag was provided.
		ok, _ := regexp.MatchString(imageRegex, image)
		if !ok {
			image = fmt.Sprintf("%s:latest", image)
		}
	}
}

func createSonarDeployment(cmd *cobra.Command, args []string) {

	k8sclient, err := k8sclient.New(kubeContext, kubeConfig)
	if err != nil {
		log.Fatal(err)
	}
}
