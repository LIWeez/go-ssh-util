package actions

import (
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/file"
	"go-ssh-util/ssh"
	"os"
	"path/filepath"
)

func RunCopyFiles() {
	selectedHost, _, err := config.ChooseAlias(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	selectedFile, isDirectoryMode, err := file.ChooseFileInteractive()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	remoteDest := file.PromptForRemoteDestination()

	userInfo := fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host)
	if isDirectoryMode {
		remoteDest = filepath.Join(remoteDest, filepath.Base(selectedFile))
		ssh.CopyUsingSCP(selectedFile, remoteDest, userInfo, selectedHost.Port, true) // Recursive copy
	} else {
		ssh.CopyUsingSCP(selectedFile, remoteDest, userInfo, selectedHost.Port, false) // Single file copy
	}
}
