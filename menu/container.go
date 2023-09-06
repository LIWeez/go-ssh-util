package menu

import (
	"fmt"
	"go-ssh-util/actions"

	"github.com/trzsz/promptui"
)

func showContainerMenu() int {
	menuPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"status", "Return to Main Menu"},
	}

	containerIndex, _, err := menuPrompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}

	return containerIndex
}

func RunContainer() {
	containerActive := true

	for containerActive {
		containerIndex := showContainerMenu()

		switch containerIndex {
		case 0:
			actions.RunStatus()

		default: // Return to main menu
			containerActive = false
		}
	}
}
