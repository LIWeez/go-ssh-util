package aws

import (
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/ssh"
	"os"
)

func RunGetVMs() {
	selectedHost, ExecutionMode, err := config.ChooseAlias(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := `aws ec2 describe-instances --region ap-southeast-2 --query "Reservations[*].Instances[*].Tags[?Key=='Name'].Value" --output text | cat`
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}
}
