package helpers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	sonartypes "github.com/glitchcrab/sonar/internal/types"
	"github.com/manifoldco/promptui"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"

	log "github.com/sirupsen/logrus"
)

func ConfirmationPrompt(resourceType, name string) bool {
	var response string

	fmt.Printf("delete %s \"%s\" [y/n]? ", resourceType, name)
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		log.Infof("not deleting %s \"%s\"", resourceType, name)
		return false
	default:
		fmt.Println("unknown response, please use 'y' or 'n':")
		return ConfirmationPrompt(resourceType, name)
	}
}

// PrintManifestYAML converts a Kubernetes object to YAML and prints it to stdout
func PrintManifestYAML(obj runtime.Object) error {
	yamlBytes, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal object to YAML: %w", err)
	}

	fmt.Println("---")
	fmt.Print(string(yamlBytes))

	return nil
}

// DisplaySelectionPrompt lists items and prompts the user to select one
func DisplaySelectionPrompt(message string, itemList []string) (selection string, err error) {
	prompt := promptui.Select{
		HideHelp:     true,
		HideSelected: true,
		Label:        message,
		Items:        itemList,
	}

	_, selection, err = prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	return selection, nil
}

// PromptForInput prompts the user for input and returns the response
func PromptForInput(promptText string) (string, error) {
	var response string

	fmt.Println(promptText)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	response = scanner.Text()

	return response, err
}

func FindSonarPods(k8sClientSet *kubernetes.Clientset, ctx context.Context, name, namespace string, searchLabels []string) ([]sonartypes.DiscoveredPod, error) {
	// Create a label selector string from the search labels.
	searchOpts := metav1.ListOptions{
		LabelSelector: strings.Join(searchLabels, ","),
	}

	// Get matching pods
	pods, err := k8sClientSet.CoreV1().Pods(namespace).List(ctx, searchOpts)
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
		if namespace == "" {
			log.Infof("no pods found with labels %s across all namespaces", strings.Join(searchLabels, ","))
		} else {
			log.Infof("no pods found with labels %s in namespace %s", strings.Join(searchLabels, ","), namespace)
		}
		os.Exit(0)
	}

	return discoveredPods, err
}

func FindSonarDeployments(k8sClientSet *kubernetes.Clientset, ctx context.Context, name, namespace string, searchLabels []string) ([]sonartypes.DiscoveredDeployment, error) {
	// Create a label selector string from the search labels.
	searchOpts := metav1.ListOptions{
		LabelSelector: strings.Join(searchLabels, ","),
	}

	// Get matching pods
	deployments, err := k8sClientSet.AppsV1().Deployments(namespace).List(ctx, searchOpts)
	if err != nil {
		log.Fatal("error listing deployments: %v", err)
	}

	var discoveredDeployments []sonartypes.DiscoveredDeployment
	for _, deploy := range deployments.Items {
		discoveredDeployments = append(discoveredDeployments, sonartypes.DiscoveredDeployment{
			Name:      deploy.Name,
			Namespace: deploy.Namespace,
		})
	}

	// Raise a clean exit if no pods found.
	if len(discoveredDeployments) == 0 {
		if namespace == "" {
			log.Infof("no deployments found with labels %s across all namespaces", strings.Join(searchLabels, ","))
		} else {
			log.Infof("no deployments found with labels %s in namespace %s", strings.Join(searchLabels, ","), namespace)
		}
		os.Exit(0)
	}

	return discoveredDeployments, err
}
