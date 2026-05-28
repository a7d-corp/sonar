package exec

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/glitchcrab/sonar/internal/app"
	"github.com/glitchcrab/sonar/internal/k8sclient"
	"github.com/glitchcrab/sonar/internal/types"
	"github.com/glitchcrab/sonar/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "exec",
		Short: "Executes a command in a Sonar debug container",
		Long: `Exec attempts to exec into a Sonar debug container using the current
(or provided) kubectl context. It searches for pods in the currently
selected namespace with the label 'owner=sonar' and prompts the user
to select a pod to exec into. The user can scope the selection by
providing a namespace via the --namespace/-n flag.

By default, the exec command will run /bin/sh in the target pod, however
any command can be provided after a '--' separator. For example:

"sonar exec -- /bin/bash" - prompts the user to select a Sonar pod
and then runs /bin/bash in the selected pod.

If the user does not provide a command via the '--' separator, they
will be prompted to enter a command after selecting a pod. If they
do not enter a command, it will default to /bin/sh.

All flags are optional.

Global flags:

Run "sonar help" in order to see flags which apply to all subcommands.`,
		Example: `
"sonar exec" - finds all Sonar pods across all namespaces.

"sonar exec --namespace kube-system" - finds all Sonar pods in
namespace 'kube-system'.`,
		RunE: runExecCommand,
	}

	return command
}

func runExecCommand(cmd *cobra.Command, args []string) error {
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

	// Create a label selector string from the search labels.
	searchOpts := metav1.ListOptions{
		LabelSelector: strings.Join(searchLabels, ","),
	}

	// Get all pods in the cluster matching the search options.
	ctx := context.TODO()
	pods, err := k8sClientSet.CoreV1().Pods(a.Globals.Namespace).List(ctx, searchOpts)
	if err != nil {
		return err
	}

	// Add all discovered pods to the list of discovered pods.
	var discoveredPods []types.DiscoveredPod
	for _, pod := range pods.Items {
		discoveredPods = append(discoveredPods, types.DiscoveredPod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    pod.Status.Phase,
		})
	}

	// Filter discovered pods to only include those in Running state, then create a list to pass to the selection prompt.
	var podList []string
	for _, pod := range discoveredPods {
		if pod.Status == corev1.PodRunning {
			podList = append(podList, fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
		}
	}

	// Raise a clean exit if no pods found.
	if len(podList) == 0 {
		if a.Globals.Namespace != "" {
			log.Infof("no running pods found with labels %s in namespace %s", strings.Join(searchLabels, ","), a.Globals.Namespace)
		} else {
			log.Infof("no running pods found with labels %s across all namespaces", strings.Join(searchLabels, ","))
		}
		return nil
	}

	// Prompt the user to select which pod to exec into.
	prompt := "Select pod to exec into"
	selectedPod, err := utils.DisplaySelectionPrompt(prompt, podList)
	if err != nil {
		return err
	}

	// Inform the user of the selected pod
	log.Infof("Selected pod: %s", selectedPod)

	// Handle passing the exec command via different methods. If the user has provided a command via the prompt, use that. If not, check if they have provided a command via the '--' separator. If not, default to /bin/sh.
	var podCommand []string

	// If the user has not provided a command via the '--' separator, prompt them to enter a command.
	if cmd.ArgsLenAtDash() < 0 {
		dynamicCommand, err := utils.PromptForInput("Enter the command to run in the pod (default: /bin/sh): ")
		if err != nil {
			log.Fatal(err)
		}

		// If the user provided a command then use it, otherwise default to /bin/sh.
		if dynamicCommand != "" {
			podCommand = strings.Split(dynamicCommand, " ")
		} else {
			podCommand = []string{"/bin/sh"}
		}
	} else {
		// Use the command provided via the '--' separator.
		podCommand = args[cmd.ArgsLenAtDash():]
	}

	log.Infof("Will run command: %s", strings.Join(podCommand, " "))

	// Trim the namespace from the selected pod.
	_, selectedPod, _ = strings.Cut(selectedPod, "/")

	// Find the target pod and namespace from the list of discovered pods based on the user-selected pod.
	var targetPod string
	var targetNamespace string
	for _, pod := range discoveredPods {
		if pod.Name == selectedPod {
			targetPod = pod.Name
			targetNamespace = pod.Namespace
		}
	}

	// Create a Kubernetes REST client for executing into the pod.
	restClient, err := k8sclient.NewRestclient(a.Globals.KubeConfig, a.Globals.KubeContext)
	if err != nil {
		return err
	}

	// Get stdin's file descriptor.
	fd := os.Stdin.Fd()

	// Exec into the pod.
	err = exec(ctx, k8sClientSet, restClient, targetPod, targetNamespace, podCommand, fd, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		return err
	}

	return nil
}
