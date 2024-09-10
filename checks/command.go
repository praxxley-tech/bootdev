package checks

import (
	"fmt"
	"strings"

	api "github.com/bootdotdev/bootdev/client"
)

// Define internal command functions
func executeLs(args []string) (string, int) {
	// Implement logic for the 'ls' command
	// For demonstration purposes, just return a static message
	return "Simulated output for 'ls'", 0
}

func executeEcho(args []string) (string, int) {
	// Implement logic for the 'echo' command
	return strings.Join(args, " "), 0
}

func executeCat(args []string) (string, int) {
	// Implement logic for the 'cat' command
	// For demonstration purposes, just return a static message
	return "Simulated output for 'cat'", 0
}

// Map allowed commands to their corresponding functions
var commandHandlers = map[string]func([]string) (string, int){
	"ls":   executeLs,
	"echo": executeEcho,
	"cat":  executeCat,
}

// Define allowed arguments for each command if needed
var allowedArgs = map[string][]string{
	"ls":   {"-l", "-a"}, // Example arguments for 'ls'
	"echo": {},           // 'echo' command can take any argument
	"cat":  {},           // 'cat' command can take any argument
}

func CLICommand(
	lesson api.Lesson,
	optionalPositionalArgs []string,
) []api.CLICommandResult {
	data := lesson.Lesson.LessonDataCLICommand.CLICommandData
	responses := make([]api.CLICommandResult, len(data.Commands))

	for i, command := range data.Commands {
		finalCommand := interpolateArgs(command.Command, optionalPositionalArgs)
		responses[i].FinalCommand = finalCommand

		cmd, args := parseCommand(finalCommand)
		if cmd == "" {
			responses[i].ExitCode = -1
			responses[i].Stdout = "Invalid command"
			continue
		}

		// Check if the command is allowed and execute it
		if handler, ok := commandHandlers[cmd]; ok {
			if !validArgs(args, allowedArgs[cmd]) {
				responses[i].ExitCode = -1
				responses[i].Stdout = "Invalid arguments"
				continue
			}

			output, exitCode := handler(args)
			responses[i].ExitCode = exitCode
			responses[i].Stdout = output
		} else {
			responses[i].ExitCode = -1
			responses[i].Stdout = "Command not allowed"
		}
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
		if len(allowedArgs) > 0 && !contains(allowedArgs, arg) {
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

// Replaces positional arguments in a command string
func interpolateArgs(rawCommand string, optionalPositionalArgs []string) string {
	for i, arg := range optionalPositionalArgs {
		rawCommand = strings.ReplaceAll(rawCommand, fmt.Sprintf("$%d", i+1), arg)
	}
	return rawCommand
}
