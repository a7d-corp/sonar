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

	"github.com/glitchcrab/sonar/internal/helpers"
	"github.com/glitchcrab/sonar/service/k8sclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	searchLabels    = []string{"owner=sonar"}
	searchNamespace string

	execCmd = &cobra.Command{
		Use:   "exec",
		Short: "Execs into a Sonar debug container",
		Long: `Exec attempts to exec into a Sonar debug container using
the current (or provided) kubectl context. It searches for pods with the
label 'owner=sonar' and uses fzf to allow the user to select a pod to
exec into. The use can further narrow the selection by providing a name
via the --name/-N flag and also the namespace via the --namespace/-n flag.

All flags are optional.

Global flags:

Run "sonar help" in order to see flags which apply to all subcommands.`,
		Example: `
"sonar exec" - finds all Sonar pods across all namespaces. If more than
one is found then user is prompted to select one.

"sonar exec --name test --namespace kube-system" - finds all Sonar pods
with the name 'test' in namespace 'kube-system'. If more than one is
found then the user is prompted to select one.`,
		Run: execIntoSonarPod,
	}
)

// DiscoveredPod represents a pod which is a candidate for execing into.
type DiscoveredPod struct {
	Name      string
	Namespace string
}

func init() {
	rootCmd.AddCommand(execCmd)
}

func execIntoSonarPod(cmd *cobra.Command, args []string) {
	log.Info("exec command is not yet implemented")

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

	var discoveredPods []DiscoveredPod
	for _, pod := range pods.Items {
		discoveredPods = append(discoveredPods, DiscoveredPod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		})
	}

	// Raise a clean exit if no pods found.
	if len(discoveredPods) == 0 {
		log.Infof("no pods found with labels %s in namespace %s", strings.Join(searchLabels, ","), searchNamespace)
		os.Exit(0)
	}

	var podList []string
	for _, pod := range discoveredPods {
		podList = append(podList, fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
	}

	// Prompt the user to select which pod to exec into.
	selectedPod, err := helpers.DisplaySelectionPrompt(podList)
	if err != nil {
		log.Fatal(err)
	}

	// Trim the namespace from the selected pod.
	_, selectedPod, _ = strings.Cut(selectedPod, "/")

	for _, targetPod := range discoveredPods {
		if targetPod.Name == selectedPod {
			//ExecIntoPod(targetPod.Namespace, targetPod.Name)
			fmt.Printf("execing into pod %s in namespace %s\n", targetPod.Name, targetPod.Namespace)
		}
	}
}
