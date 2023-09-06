package menu

import (
	"fmt"
	"github.com/trzsz/promptui"
	"go-ssh-util/actions"
)

func showNetworkMenu() int {
	menuPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"ping", "Return to Main Menu"},
	}

	networkIndex, _, err := menuPrompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}

	return networkIndex
}

func RunNetwork() {
	networkActive := true

	for networkActive {
		networkIndex := showNetworkMenu()

		switch networkIndex {
		case 0:
			actions.RunPing()
		default: // Return to main menu
			networkActive = false
		}
	}
}
