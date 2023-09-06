package menu

import (
	"fmt"
	"go-ssh-util/actions"

	"github.com/trzsz/promptui"
)

func showMainMenu() int {
	menuPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"\U0001F4C1 Copy File/Directory (SCP)", "\U0001F5A5  Check Systeminfo", "\U0001F4E1 Check Network", "\U0001F433 Check Container", "\U0001F511 SSH", "\U0001F30E Cloud", "Quit"},
		Size:  10,
	}

	scpIndex, _, err := menuPrompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}

	return scpIndex
}

func RunMainLoop() {
	for {
		scpIndex := showMainMenu()

		switch scpIndex {
		case 0: // Copy File/Directory (SCP)
			actions.RunCopyFiles()
		case 1: // systeminfo
			RunSystemInfo()
		case 2: // network
			RunNetwork()
		case 3:
			RunContainer()
		case 4:
			actions.RunSSH()
		case 5:
			RunCloud()
		case 6:
			fmt.Println("Goodbye!")
			return
		}

	}
}
