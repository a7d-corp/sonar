package utils

import (
	"context"
	"os"
	"strings"

	sonartypes "github.com/glitchcrab/sonar/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	log "github.com/sirupsen/logrus"
)

// FindSonarPods searches for pods matching the provided labels and returns a list of discovered pods.
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
