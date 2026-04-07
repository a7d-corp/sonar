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

// FindSonarDeployments searches for Kubernetes deployments matching the provided labels and returns a list of discovered deployments.
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
