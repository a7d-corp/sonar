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
package ls

import (
	"context"
	"fmt"
	"strings"

	"github.com/glitchcrab/sonar/internal/app"
	"github.com/glitchcrab/sonar/internal/types"
	"github.com/glitchcrab/sonar/service/k8sclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list", "discover", "lookup"},
		Short:   "Lists all Sonar debug containers",
		Long: `ls attempts to discover all debug containers in the cluster
which were created by Sonar. It searches for pods with the label
'owner=sonar' and lists them for the user. The user can further
narrow the selection by providing a name via the --name/-N flag
and also the namespace via the --namespace/-n flag.

All flags are optional.

Global flags:

Run "sonar help" in order to see flags which apply to all subcommands.`,
		Example: `
"sonar ls" - finds all Sonar pods across all namespaces.

"sonar ls --name test --namespace kube-system" - finds all Sonar pods
with the name 'test' in namespace 'kube-system'.`,
		RunE: runLsCommand,
	}

	return command
}

func runLsCommand(cmd *cobra.Command, args []string) error {
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

	// Assemble lookup options

	// Labels used to match Sonar containers.
	searchLabels := []string{"owner=sonar"}

	// Add the provided name to the search labels if it is not empty.
	if a.Globals.Name != "" {
		// add the name to the search labels - we use the full name value as
		// this is what the pod is actually labelled with.
		searchLabels = append(searchLabels, fmt.Sprintf("name=%s", a.Globals.Name))
	}

	// Determine whether to scope the search to a specific namespace or across the whole cluster.
	var searchNamespace string
	if a.Globals.Namespace != "" {
		searchNamespace = a.Globals.Namespace
	}

	// Create a label selector string from the search labels.
	searchOpts := metav1.ListOptions{
		LabelSelector: strings.Join(searchLabels, ","),
	}

	// Get all pods in the cluster matching the search options.
	ctx := context.TODO()
	pods, err := k8sClientSet.CoreV1().Pods(searchNamespace).List(ctx, searchOpts)
	if err != nil {
		return err
	}

	// Add all discovered pods to the list of discovered pods.
	discoveredPods := []types.DiscoveredPod{}
	for _, pod := range pods.Items {
		discoveredPods = append(discoveredPods, types.DiscoveredPod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    pod.Status.Phase,
		})
	}

	// Raise a clean exit if no pods found.
	if len(discoveredPods) == 0 {
		if a.Globals.Namespace != "" {
			log.Infof("no pods found with labels %s in namespace %s", strings.Join(searchLabels, ","), a.Globals.Namespace)
		} else {
			log.Infof("no pods found with labels %s across all namespaces", strings.Join(searchLabels, ","))
		}
		return nil
	}

	// Print all discovered pods.
	for _, pod := range discoveredPods {
		log.Infof("found %s in namespace %s (status: %s)", pod.Name, pod.Namespace, pod.Status)
	}

	return nil
}
