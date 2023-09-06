package actions

import (
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/ssh"
	"os"
)

func RunStatus() {
	selectedHost, ExecutionMode, err := config.ChooseAlias(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := "docker ps"
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}

}
