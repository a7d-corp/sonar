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
	"os"
	"strings"

	sonartypes "github.com/glitchcrab/sonar/internal/types"
	"github.com/glitchcrab/sonar/service/k8sclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	lsCmd = &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list", "find", "search", "discover", "lookup"},
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
		Run: lsCommand,
	}
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

func lsCommand(cmd *cobra.Command, args []string) {
	var searchLabels = []string{"owner=sonar"}
	var searchNamespace string

	// Create a clientset to interact with the cluster (only if not in dry-run mode).
	var k8sClientSet *kubernetes.Clientset
	var err error

	k8sClientSet, err = k8sclient.New(kubeContext, kubeConfig)
	if err != nil {
		log.Fatal(err) // TODO: better logging
	}

	ctx := context.TODO()

	// Assemble lookup options

	/// Set whether we search all namespaces or scope to a specific namespace.
	if rootCmd.PersistentFlags().Lookup("namespace").Changed {
		// use the provided namespace
		searchNamespace = namespace
	} else {
		// search all namespaces
		searchNamespace = ""
	}

	if rootCmd.PersistentFlags().Lookup("name").Changed {
		// add the name to the search labels - we use the full name value as
		// this is what the pod is actually labelled with.
		searchLabels = append(searchLabels, fmt.Sprintf("name=%s", fullName))
	}

	// Create a label selector string from the search labels.
	searchOpts := metav1.ListOptions{
		LabelSelector: strings.Join(searchLabels, ","),
	}

	pods, err := k8sClientSet.CoreV1().Pods(searchNamespace).List(ctx, searchOpts)
	if err != nil {
		log.Fatal("error listing pods: %v", err)
	}

	var discoveredPods []sonartypes.DiscoveredPod
	for _, pod := range pods.Items {
		discoveredPods = append(discoveredPods, sonartypes.DiscoveredPod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    pod.Status.Phase,
		})
	}

	// Raise a clean exit if no pods found.
	if len(discoveredPods) == 0 {
		if searchNamespace == "" {
			log.Infof("no pods found with labels %s across all namespaces", strings.Join(searchLabels, ","))
		} else {
			log.Infof("no pods found with labels %s in namespace %s", strings.Join(searchLabels, ","), searchNamespace)
		}
		os.Exit(0)
	}

	for _, pod := range discoveredPods {
		log.Infof("found %s in namespace %s (status: %s)", pod.Name, pod.Namespace, pod.Status)
	}
}
