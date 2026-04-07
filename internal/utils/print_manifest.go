package utils

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

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
