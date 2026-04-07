package utils

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

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
