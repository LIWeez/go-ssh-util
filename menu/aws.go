package menu

import (
	"fmt"
	"go-ssh-util/actions/aws"

	"github.com/trzsz/promptui"
)

func showAWSMenu() int {
	menuPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"EC2", "Return to Main Menu"},
	}

	AWSIndex, _, err := menuPrompt.Run()
	if err != nil {
		fmt.Println("Error", err)
		return -1
	}
	return AWSIndex
}

func RunAWS() {
	AWSActive := true

	for AWSActive {
		AWSIndex := showAWSMenu()

		switch AWSIndex {
		case 0:
			RunEC2()
		default:
			AWSActive = false
		}
	}
}

func showEC2Menu() int {
	menuPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"Get VMs"},
	}

	AWSIndex, _, err := menuPrompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}

	return AWSIndex
}

func RunEC2() {
	AWSActive := true

	for AWSActive {
		AWSIndex := showEC2Menu()

		switch AWSIndex {
		case 0:
			aws.RunGetVMs()
		default:
			AWSActive = false
		}
	}
}
