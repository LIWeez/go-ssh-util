package actions

import (
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/ssh"
	"os"
)

func RunSSH() {
	selectedHost, _, err := config.ChooseAlias(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	ssh.SSHToRemoteHostWithKey(selectedHost)
}
