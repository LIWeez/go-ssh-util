package actions

import (
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/ssh"
	"os"
)

func RunCheckCPU() {
	selectedHost, ExecutionMode, err := config.ChooseAlias(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := "lscpu"
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}
}

func RunCheckMemory() {
	selectedHost, ExecutionMode, err := config.ChooseAlias(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := "free -h"
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}
}

func RunCheckDisk() {
	selectedHost, ExecutionMode, err := config.ChooseAlias(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := "df -h"
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}
}

func RunCheckSwap() {
	selectedHost, ExecutionMode, err := config.ChooseAlias(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := "cat /proc/swaps"
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}
}

func RunCheckNetwork() {
	selectedHost, ExecutionMode, err := config.ChooseAlias(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := "ifconfig"
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}
}
