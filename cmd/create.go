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
	"fmt"
	"regexp"

	"github.com/glitchcrab/sonar/internal/sonarconfig"
	"github.com/glitchcrab/sonar/service/k8sclient"
	"github.com/glitchcrab/sonar/service/k8sresource"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	imageRegex = "^[a-z0-9/.-]*[:][a-z0-9.-]*$"
)

var (
	image             string
	networkPolicy     bool
	podArgs           string
	podCommand        string
	podSecurityPolicy bool
	podUser           int64
	privileged        bool

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "create applies a debug deployment to a Kubernetes cluster",
		Long: `create will attempt to create a debugging deployment and all supporting
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
	createCmd.Flags().BoolVar(&networkPolicy, "networkpolicy", false, "create NetworkPolicy (default \"false\")")
	createCmd.Flags().StringVarP(&podArgs, "pod-args", "a", "24h", "args to pass to pod command")
	createCmd.Flags().StringVarP(&podCommand, "pod-command", "c", "sleep", "pod command (aka image entrypoint)")
	createCmd.Flags().BoolVar(&podSecurityPolicy, "podsecuritypolicy", false, "create PodSecurityPolicy (default \"false\")")
	createCmd.Flags().Int64VarP(&podUser, "pod-userid", "u", 1000, "userID to run the pod as")
	createCmd.Flags().BoolVar(&privileged, "privileged", false, "run the container as root (assumes userID of 0) (default \"false\")")
}

func validateFlags() {
	// If the user has provided an image name then validate that it looks sane.
	if createCmd.Flags().Lookup("image").Changed {
		// Validate image to see if a tag has been provided; if not then
		// use :latest. Does not validate full image name, just whether a
		// tag was provided.
		ok, _ := regexp.MatchString(imageRegex, image)
		if !ok {
			image = fmt.Sprintf("%s:latest", image)
		}
	}
}

func createSonarDeployment(cmd *cobra.Command, args []string) {
	// Create a SonarConfig
	sonarConfig := sonarconfig.SonarConfig{
		Image:             image,
		Labels:            labels,
		Name:              name,
		Namespace:         namespace,
		NetworkPolicy:     networkPolicy,
		PodArgs:           podArgs,
		PodCommand:        podCommand,
		PodSecurityPolicy: podSecurityPolicy,
		PodUser:           podUser,
		Privileged:        privileged,
	}
	fmt.Println(sonarConfig)
	// Create a clientset to interact with the cluster.
	k8sClientSet, err := k8sclient.New(kubeContext, kubeConfig)
	if err != nil {
		log.Fatal(err) // TODO: better logging
	}

	// Create a context
	ctx := context.TODO()

	{
		// Create a ServiceAccount
		err := k8sresource.NewServiceAccount(k8sClientSet, ctx, sonarConfig)
		// Handle the response
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			log.Warnf("serviceaccount \"%s/%s\" already exists\n", sonarConfig.Namespace, sonarConfig.Name)
		} else if err != nil {
			log.Warnf("serviceaccount \"%s/%s\" was not created: %w\n", sonarConfig.Namespace, sonarConfig.Name, err)
		} else {
			log.Infof("serviceaccount \"%s/%s\" created\n", sonarConfig.Namespace, sonarConfig.Name)
		}
	}

	if sonarConfig.PodSecurityPolicy {
		{
			// Create a PodSecurityPolicy
			err := k8sresource.NewPodSecurityPolicy(k8sClientSet, ctx, sonarConfig)
			// Handle the response
			if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
				log.Warnf("podsecuritypolicy \"%s\" already exists\n", sonarConfig.Name)
			} else if err != nil {
				log.Warnf("podsecuritypolicy \"%s\" was not created: %w\n", sonarConfig.Name, err)
			} else {
				log.Infof("podsecuritypolicy \"%s\" created\n", sonarConfig.Name)
			}
		}

		{
			// Create a ClusterRole
			err := k8sresource.NewClusterRole(k8sClientSet, ctx, sonarConfig)
			// Handle the response
			if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
				log.Warnf("clusterrole \"%s\" already exists\n", sonarConfig.Name)
			} else if err != nil {
				log.Warnf("clusterrole \"%s\" was not created: %w\n", sonarConfig.Name, err)
			} else {
				log.Infof("clusterrole \"%s\" created\n", sonarConfig.Name)
			}
		}

		{
			// Create a ClusterRoleBinding
			err := k8sresource.NewClusterRoleBinding(k8sClientSet, ctx, sonarConfig)
			// Handle the response
			if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
				log.Warnf("clusterrolebinding \"%s\" already exists\n", sonarConfig.Name)
			} else if err != nil {
				log.Warnf("clusterrolebinding \"%s\" was not created: %w\n", sonarConfig.Name, err)
			} else {
				log.Infof("clusterrolebinding \"%s\" created\n", sonarConfig.Name)
			}
		}
	}

	if sonarConfig.NetworkPolicy {
		// Create a NetworkPolicy
		err := k8sresource.NewNetworkPolicy(k8sClientSet, ctx, sonarConfig)
		// Handle the response
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			log.Warnf("networkpolicy \"%s\" already exists\n", sonarConfig.Name)
		} else if err != nil {
			log.Warnf("networkpolicy \"%s\" was not created: %w\n", sonarConfig.Name, err)
		} else {
			log.Infof("networkpolicy \"%s\" created\n", sonarConfig.Name)
		}
	}

	{
		// Create a Deployment
		err := k8sresource.NewDeployment(k8sClientSet, ctx, sonarConfig)
		// Handle the response
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonAlreadyExists {
			log.Warnf("deployment \"%s/%s\" already exists\n", sonarConfig.Namespace, sonarConfig.Name)
		} else if err != nil {
			log.Warnf("deployment \"%s/%s\" was not created: %w\n", sonarConfig.Namespace, sonarConfig.Name, err)
		} else {
			log.Infof("deployment \"%s/%s\" created\n", sonarConfig.Namespace, sonarConfig.Name)
		}
	}
}
