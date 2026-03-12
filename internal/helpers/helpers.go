package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
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

// DisplaySelectionPrompt lists items and prompts the user to select one
func DisplaySelectionPrompt(itemList []string) (selection string, err error) {
	prompt := promptui.Select{
		Label: "Select an item",
		Items: itemList,
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
