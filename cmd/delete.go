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
	"github.com/glitchcrab/sonar/internal/helpers"
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

Skips conformation prompt and deletes all resources created by Sonar.`,
		Example: `
"sonar delete" - prompts the user to select a Sonar deployment from a
list of all matching deployments in a cluster.

"sonar delete --name test" - deletes all resources named 'test'.
in namespace 'kube-system' named 'sonar-test'.

NOTE: passing the --namespace flag will scope the search for deployments
to a specific namespace.`,
		Run: deleteSonarDeployment,
	}
)

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&force, "force", "f", false, "skip all confirmation prompts when deleting (default \"false\")")
}

func deleteSonarDeployment(cmd *cobra.Command, args []string) {
	var searchLabels = []string{"owner=sonar"}
	var searchNamespace string
	var selectedDeploy string

	// Create a clientset to interact with the cluster.
	k8sClientSet, err := k8sclient.New(kubeContext, kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx := context.TODO()

	// Check if the user provided the --name flag, if so we skip the interactive lookup.
	var skipInteractiveLookup bool
	if rootCmd.PersistentFlags().Lookup("name").Changed {
		skipInteractiveLookup = true
		log.Info("name flag was set, skipping interactive selection")
	}

	/// Set whether we search all namespaces or scope to a specific namespace.
	if rootCmd.PersistentFlags().Lookup("namespace").Changed {
		// use the provided namespace
		searchNamespace = namespace
	} else {
		// search all namespaces
		searchNamespace = ""
	}

	discoveredDeployments, err := helpers.FindSonarDeployments(k8sClientSet, ctx, name, searchNamespace, searchLabels)

	if !skipInteractiveLookup {
		// Build a list of deployments to pass to the selection prompt.
		var deployList []string
		for _, deploy := range discoveredDeployments {
			deployList = append(deployList, fmt.Sprintf("%s/%s", deploy.Namespace, deploy.Name))
		}

		// Prompt the user to select which deployment to delete.
		prompt := "Select deployment to delete"
		selectedDeploy, err = helpers.DisplaySelectionPrompt(prompt, deployList)
		if err != nil {
			log.Fatal(err)
		}

		// Trim the namespace from the selected deployment.
		_, selectedDeploy, _ = strings.Cut(selectedDeploy, "/")
	} else {
		// If the user provided a name, we use that as the selected deployment.
		selectedDeploy = fullName
	}

	// Inform the user of the selected pod
	log.Infof("Deployment to be deleted: %s", selectedDeploy)

	// Find the namespace of the victim deployment.
	var selectedDeployNamespace string
	for _, d := range discoveredDeployments {
		if d.Name == selectedDeploy {
			selectedDeployNamespace = d.Namespace
			break
		}
	}

	// Create a SonarConfig and populate it with enough variables for deletion.
	sonarConfig := clientconfigs.SonarConfig{
		Labels:    labels,
		Name:      selectedDeploy,
		Namespace: selectedDeployNamespace,
	}

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
