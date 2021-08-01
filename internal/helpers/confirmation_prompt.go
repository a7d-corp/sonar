package helpers

import (
	"fmt"
	"strings"

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
		return false
	default:
		fmt.Println("unknown response, please use 'y' or 'n':")
		return ConfirmationPrompt(resourceType, name)
	}
}
