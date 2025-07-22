package helpers

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
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
