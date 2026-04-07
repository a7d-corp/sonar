package utils

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

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
