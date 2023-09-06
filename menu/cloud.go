package menu

import (
	"fmt"

	"github.com/trzsz/promptui"
)

func showCloudMenu() int {
	menuPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"GCP", "AWS", "Return to Main Menu"},
	}

	CloudIndex, _, err := menuPrompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}

	return CloudIndex
}

func RunCloud() {
	CloudActive := true

	for CloudActive {
		CloudIndex := showCloudMenu()

		switch CloudIndex {
		case 0:
			RunGCP()

		case 1:
			RunAWS()

		default: // Return to main menu
			CloudActive = false
		}
	}
}
