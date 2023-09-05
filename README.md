# Go SSH Util

Go SSH Util is a versatile command-line tool that empowers you to securely manage remote systems using SSH. It provides various features to facilitate secure file transfers, remote command execution, and system monitoring, all within an intuitive and user-friendly interface.

## Features

- **Secure Copy (SCP):** Transfer files and directories between your local machine and remote servers securely using SCP protocol.
- **Remote Command Execution:** Execute commands on remote servers securely using SSH, making it easy to manage and monitor remote systems.
- **System Monitoring:** Monitor system resources like memory usage, disk space, and swap space on remote servers.

## Installation

Make sure you have Go installed on your system. Then, install the tool using the following command:

```sh
go install github.com/sz9751210/go-ssh-util@latest
```

## Prerequisites

- Go (Golang) installed on your system.
- Basic familiarity with the command-line interface.

## How to Use

1. Clone or download the repository to your local machine.

2. Open a terminal and navigate to the directory containing the cloned repository.

3. Run the following command to build and execute the program:
   
   ```bash
   go run cmd/main.go
   ```
    The program will guide you through the process of selecting a file, configuring the remote host details, and initiating the file copy.

4. Follow the on-screen prompts to navigate through directories, choose a file, and configure the copy settings.

5. Once the configuration is complete, the program will use SCP to copy the selected file to the specified remote host.

## Contributing
Contributions are welcome! If you find a bug or have an idea for improvement, please open an issue or submit a pull request on the GitHub repository.

## License
This project is licensed under the [MIT License](LICENSE).