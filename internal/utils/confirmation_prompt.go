package utils

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ConfirmationPrompt prompts the user for confirmation before deleting a resource. It returns true if the user confirms, and false otherwise.
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
