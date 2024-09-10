package checks

import (
	"os/exec"
	"strings"

	api "github.com/bootdotdev/bootdev/client"
)

// Whitelist commands with predefined arguments
var allowedCommands = map[string][]string{
	"ls":   {"-l", "-a"}, // Example arguments for `ls`
	"echo": {},           // `echo` command can take any argument
	"cat":  {},           // `cat` command can take any argument
	// Add more allowed commands with arguments as needed
}

func CLICommand(
	lesson api.Lesson,
	optionalPositionalArgs []string,
) []api.CLICommandResult {
	data := lesson.Lesson.LessonDataCLICommand.CLICommandData
	responses := make([]api.CLICommandResult, len(data.Commands))

	for i, command := range data.Commands {
		// Use predefined commands and arguments only
		cmd, args := parseCommand(command.Command)
		if cmd == "" {
			responses[i].ExitCode = -1
			responses[i].Stdout = "Invalid command"
			continue
		}

		// Check whitelist
		if allowedArgs, ok := allowedCommands[cmd]; ok {
			if !validArgs(args, allowedArgs) {
				responses[i].ExitCode = -1
				responses[i].Stdout = "Invalid arguments"
				continue
			}
		} else {
			responses[i].ExitCode = -1
			responses[i].Stdout = "Command not allowed"
			continue
		}

		// Execute the command securely
		execCmd := exec.Command(cmd, args...)

		// Capture output
		b, err := execCmd.CombinedOutput()
		if ee, ok := err.(*exec.ExitError); ok {
			responses[i].ExitCode = ee.ExitCode()
		} else if err != nil {
			responses[i].ExitCode = -2
		} else {
			responses[i].ExitCode = 0
		}

		// Store output
		responses[i].Stdout = strings.TrimRight(string(b), " \n\t\r")
	}

	return responses
}

// Parses a command string into command and arguments
func parseCommand(command string) (string, []string) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

// Validates arguments against allowed arguments for a command
func validArgs(args []string, allowedArgs []string) bool {
	for _, arg := range args {
		if len(arg) == 0 {
			return false
		}
		if !contains(allowedArgs, arg) {
			return false
		}
	}
	return true
}

// Checks if a slice contains a specific element
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}
