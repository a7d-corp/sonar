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
	"strings"

	"github.com/glitchcrab/sonar/internal/clientconfigs"
	"github.com/glitchcrab/sonar/service/k8sclient"
	"github.com/glitchcrab/sonar/service/k8sresource"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	force bool

	deleteCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"destroy"},
		Short:   "Delete destroys all Sonar resources",
		Long: `Delete will attempt to remove all resources deployed to a cluster
by Sonor in the provided kubectl context (or the current context if
none is provided).

All flags are optional, however if the deployment was configured when
it was initially deployed then a combination of flags will be required
in order to ensure that Sonar can find the resources.

Global flags:

Run "sonar help" in order to see flags which apply to all subcommands.

Flags:

--force (default: false)

Skips all interaction and deletes all resources created by Sonar.`,
		Example: `
"sonar delete" - deletes all resources which match the defaults. This
will result in deleting all resources in namespace 'default' which are
named 'sonar-debug'.

"sonar delete --name test --namespace kube-system" - deletes all
resources in namespace 'kube-system' named 'sonar-test'.`,
		Run: deleteSonarDeployment,
	}
)

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&force, "force", "f", false, "skip all confirmation prompts when deleting (default \"false\")")
}

func deleteSonarDeployment(cmd *cobra.Command, args []string) {
	// Create a SonarConfig and populate it with enough variables for deletion.
	sonarConfig := clientconfigs.SonarConfig{
		Labels:    labels,
		Name:      fullName,
		Namespace: namespace,
	}

	// Create a clientset to interact with the cluster.
	k8sClientSet, err := k8sclient.New(kubeContext, kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx := context.TODO()

	// Initialise an empty map to report deleted resources to user
	// at the end.
	deletedResources := []string{}

	if force {
		log.Info("force was set, not asking for confirmation before deleting resources")
	}

	{
		// Delete the deployment
		err := k8sresource.DeleteDeployment(k8sClientSet, ctx, sonarConfig, force)
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
			log.Info("no matching deployment found; skipping deletion")
		} else if err != nil {
			log.Warnf("deployment \"%s/%s\" failed deletion: %w", sonarConfig.Namespace, sonarConfig.Name, err)
		} else {
			log.Infof("deleting deployment")
			deletedResources = append(deletedResources, "deployment")
		}
	}

	// Convert sonarConfig.Labels into a format suitable for use as a LabelSelector.
	var labelSlice = []string{}
	for k, v := range labels {
		labelSlice = append(labelSlice, fmt.Sprintf("%s=%s", k, v))
	}

	// Filter resources by Sonar labels.
	listOpts := metav1.ListOptions{
		LabelSelector: strings.Join(labelSlice, ","),
	}

	{
		// Get NetworkPolicies and see if a match is found
		inClusterNps := []networkingv1.NetworkPolicy{}
		nps, err := k8sClientSet.NetworkingV1().NetworkPolicies(sonarConfig.Namespace).List(ctx, listOpts)
		if err != nil {
			log.Warnf("%w", err) // TODO: improve error logging here
		}
		inClusterNps = append(inClusterNps, nps.Items...)

		// Range over discovered NetworkPolicies and see if any match.
		for _, np := range inClusterNps {
			if strings.HasPrefix(np.Name, sonarConfig.Name) {
				// Delete the NetworkPolicy
				err := k8sresource.DeleteNetworkPolicy(k8sClientSet, ctx, sonarConfig, force)
				if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
					log.Info("no matching podsecuritypolicy found; skipping deletion")
				} else if err != nil {
					log.Warnf("networkpolicy \"%s\" failed deletion: %w", sonarConfig.Name, err)
				} else {
					log.Infof("deleting networkpolicy")
					deletedResources = append(deletedResources, "networkpolicy")
				}
			}
		}
	}

	{
		// Delete the ServiceAccount
		err := k8sresource.DeleteServiceAccount(k8sClientSet, ctx, sonarConfig, force)
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
			log.Info("no matching serviceaccount found; skipping deletion")
		} else if err != nil {
			log.Warnf("serviceaccount \"%s/%s\" failed deletion: %w", sonarConfig.Namespace, sonarConfig.Name, err)
		} else {
			log.Infof("deleting serviceaccount")
			deletedResources = append(deletedResources, "serviceaccount")
		}
	}

	if len(deletedResources) > 0 {
		log.Infof("resources deleted: %s", strings.Join(deletedResources, ", "))
	} else {
		log.Info("no resources were deleted")
	}
}
