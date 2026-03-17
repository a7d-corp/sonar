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
	"io"
	"os"
	"strings"

	"github.com/glitchcrab/sonar/internal/helpers"
	sonartypes "github.com/glitchcrab/sonar/internal/types"
	"github.com/glitchcrab/sonar/service/k8sclient"
	"github.com/moby/term"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	execCmd = &cobra.Command{
		Use:   "exec",
		Short: "Execs into a Sonar debug container",
		Long: `Exec attempts to exec into a Sonar debug container using the current
(or provided) kubectl context. It searches for pods with the label
'owner=sonar' and prompts the user to select a pod to exec into. The
user can further narrow the selection by providing a name via the
--name/-N flag and also the namespace via the --namespace/-n flag.

By default, the exec command will run /bin/sh in the target pod, however
any command can be provided after a '--' separator. For example:

"sonar exec -- /bin/bash" - prompts the user to select a Sonar pod
and then runs /bin/bash in the selected pod.

All flags are optional.

Global flags:

Run "sonar help" in order to see flags which apply to all subcommands.`,
		Example: `
"sonar exec" - finds all Sonar pods across all namespaces.

"sonar exec --name test --namespace kube-system" - finds all Sonar pods
with the name 'test' in namespace 'kube-system'.`,
		Run: execCommand,
	}
)

func init() {
	rootCmd.AddCommand(execCmd)
}

func execCommand(cmd *cobra.Command, args []string) {
	var searchLabels = []string{"owner=sonar"}
	var searchNamespace string
	var targetPod string
	var targetNamespace string

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

	var podList []string
	for _, pod := range discoveredPods {
		podList = append(podList, fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
	}

	// Prompt the user to select which pod to exec into.
	selectedPod, err := helpers.DisplaySelectionPrompt(podList)
	if err != nil {
		log.Fatal(err)
	}

	// Handle passing the exec command via different methods. If the user has provided a command via the prompt, use that. If not, check if they have provided a command via the '--' separator. If not, default to /bin/sh.
	var podCommand []string

	// If the user has not provided a command via the '--' separator, prompt them to enter a command.
	if cmd.ArgsLenAtDash() < 0 {
		dynamicCommand, err := helpers.PromptForInput("Enter the command to run in the pod (or leave it blank)")
		if err != nil {
			log.Fatal(err)
		}

		// If the user provided a command then use it, otherwise default to /bin/sh.
		if dynamicCommand == "" {
			podCommand = []string{"/bin/sh"}
		} else {
			podCommand = strings.Split(dynamicCommand, " ")
		}
	} else {
		// Use the command provided via the '--' separator.
		podCommand = args[cmd.ArgsLenAtDash():]
	}

	log.Infof("Will run command: %s", strings.Join(podCommand, " "))

	// Trim the namespace from the selected pod.
	_, selectedPod, _ = strings.Cut(selectedPod, "/")

	for _, pod := range discoveredPods {
		if pod.Name == selectedPod {
			targetPod = pod.Name
			targetNamespace = pod.Namespace
		}
	}

	restClient, err := k8sclient.NewRestclient(kubeConfig, kubeContext)
	if err != nil {
		log.Fatal(err)
	}

	fd := os.Stdin.Fd()
	err = execIntoPod(ctx, k8sClientSet, restClient, targetPod, targetNamespace, podCommand, fd, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		log.Fatal(err)
	}
}

// execIntoPod execs into a pod and gets a shell
func execIntoPod(ctx context.Context, k8sClientSet *kubernetes.Clientset, restClient *restclient.Config, targetPod, targetNamespace string, podCommand []string, fd uintptr, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	request := k8sClientSet.CoreV1().RESTClient().Post().Resource("pods").Name(targetPod).Namespace(targetNamespace).SubResource("exec")
	options := &corev1.PodExecOptions{
		Command: podCommand,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	request.VersionedParams(
		options,
		scheme.ParameterCodec,
	)

	executor, err := remotecommand.NewSPDYExecutor(restClient, "POST", request.URL())
	if err != nil {
		return err
	}

	streamOpts := remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}

	log.Infof("Connecting to pod %s in namespace %s, use Ctrl+d to exit\n\n", targetPod, targetNamespace)

	var previousState *term.State
	previousState, err = term.SetRawTerminal(fd)
	if err != nil {
		log.Fatal(err)
	}

	defer term.RestoreTerminal(fd, previousState)
	if err != nil {
		log.Fatal(err)
	}

	err = executor.StreamWithContext(ctx, streamOpts)
	if err != nil {
		return err
	}

	return nil
}
