package menu

import (
	"fmt"
	"github.com/trzsz/promptui"
	"go-ssh-util/actions"
)

func showSystemInfoMenu() int {
	menuPrompt := promptui.Select{
		Label: "Choose an operation:",
		Items: []string{"Check CPU", "Check Memory", "Check Disk", "Check Swap", "Check Network", "Return to Main Menu"},
	}

	systeminfoIndex, _, err := menuPrompt.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}

	return systeminfoIndex
}

func RunSystemInfo() {
	systemInfoActive := true

	for systemInfoActive {
		systeminfoIndex := showSystemInfoMenu()

		switch systeminfoIndex {
		case 0: // Check CPU
			actions.RunCheckCPU()
		case 1: // Check Memory
			actions.RunCheckMemory()
		case 2: // Check Disk
			actions.RunCheckDisk()
		case 3: // Check Swap
			actions.RunCheckSwap()
		case 4: // Check Network
			actions.RunCheckNetwork()
		default: // Return to main menu
			systemInfoActive = false
		}
	}
}
