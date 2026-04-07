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
package destroy

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/glitchcrab/sonar/internal/app"
	"github.com/glitchcrab/sonar/internal/config"
	"github.com/glitchcrab/sonar/internal/helpers"
	"github.com/glitchcrab/sonar/internal/k8sclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	force bool
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "destroy",
		Aliases: []string{"delete", "remove", "rm"},
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
		RunE: runDeleteCommand,
	}

	command.Flags().BoolVarP(&force, "force", "f", false, "skip all confirmation prompts when deleting")

	return command
}

func runDeleteCommand(cmd *cobra.Command, args []string) error {
	// Get the App instance from the command context
	a, err := app.GetApp(cmd)
	if err != nil {
		return err
	}

	// Create a Kubernetes clientset.
	k8sClientSet, err := k8sclient.New(a.Globals.KubeContext, a.Globals.KubeConfig)
	if err != nil {
		return err
	}

	// Check if the user provided the --name flag, if so we skip the interactive lookup.
	var skipInteractiveLookup bool
	if a.Globals.Name != "" {
		skipInteractiveLookup = true
		log.Info("name flag was set, skipping interactive selection")
	}

	/// Set whether we search all namespaces or scope to a specific namespace.
	var searchNamespace string
	if a.Globals.Namespace != "" {
		// use the provided namespace
		searchNamespace = a.Globals.Namespace
	} else {
		// search all namespaces
		searchNamespace = ""
	}

	// Labels used to match Sonar resources.
	searchLabels := []string{"owner=sonar"}

	// Add the provided name to the search labels if it is not empty.
	if a.Globals.Name != "" {
		// add the name to the search labels - we use the full name value as
		// this is what the pod is actually labelled with.
		searchLabels = append(searchLabels, fmt.Sprintf("name=%s", a.Globals.Name))
	}

	// Create a context
	ctx := context.TODO()

	// Find all Sonar deployments.
	discoveredDeployments, err := helpers.FindSonarDeployments(k8sClientSet, ctx, a.Globals.Name, searchNamespace, searchLabels)

	var selectedDeploy string
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
		selectedDeploy = a.Globals.FullName
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

	// Instantiate a DeleteConfig struct.
	opts := config.DeleteConfig{
		SearchLabels: searchLabels,
		Name:         selectedDeploy,
		Namespace:    selectedDeployNamespace,
	}

	if force {
		log.Info("force was set, not asking for confirmation before deleting resources")
	}

	// Collect any errors.
	var errs []error

	// Initialise an empty map to report deleted resources to user
	// at the end.
	deletedResources := []string{}

	// Convert opts.Labels into a format suitable for use as a LabelSelector.
	//var labelSlice = []string{}
	//for k, v := range opts.SearchLabels {
	//	labelSlice = append(labelSlice, fmt.Sprintf("%s=%s", k, v))
	//}

	// Filter resources by Sonar labels.
	listOpts := metav1.ListOptions{
		LabelSelector: strings.Join(opts.SearchLabels, ","),
	}

	// Delete the deployment
	delDeploy, deployErr := deleteDeployment(k8sClientSet, ctx, opts, force)
	if deployErr != nil {
		errs = append(errs, deployErr)
	}
	// If the deployment was deleted successfully, add it to the list of deleted resources.
	if len(delDeploy) > 0 {
		deletedResources = append(deletedResources, delDeploy)
	}

	// Get NetworkPolicies and see if a match is found
	inClusterNps := []networkingv1.NetworkPolicy{}
	nps, err := k8sClientSet.NetworkingV1().NetworkPolicies(opts.Namespace).List(ctx, listOpts)
	if err != nil {
		log.Warnf("%w", err)
	}
	inClusterNps = append(inClusterNps, nps.Items...)

	// Range over discovered NetworkPolicies and see if any match.
	for _, np := range inClusterNps {
		if strings.HasPrefix(np.Name, opts.Name) {
			// Delete the NetworkPolicy
			delNP, npErr := deleteNetworkPolicy(k8sClientSet, ctx, opts, force)
			if npErr != nil {
				errs = append(errs, npErr)
			}
			// If the NetworkPolicy was deleted successfully, add it to the list of deleted resources.
			if len(delNP) > 0 {
				deletedResources = append(deletedResources, delNP)
			}
		}
	}

	// Delete the ServiceAccount
	delSA, saErr := deleteServiceAccount(k8sClientSet, ctx, opts, force)
	if saErr != nil {
		errs = append(errs, saErr)
	}
	// If the ServiceAccount was deleted successfully, add it to the list of deleted resources.
	if len(delSA) > 0 {
		deletedResources = append(deletedResources, delSA)
	}

	// If there were any validation errors, return them as a single error.
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	if len(deletedResources) > 0 {
		log.Infof("resources deleted: %s", strings.Join(deletedResources, ", "))
	} else {
		log.Info("no resources were deleted")
	}

	return nil
}
