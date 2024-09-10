package checks

import (
	"fmt"
	"os/exec"
	"strings"
	"unicode"

	api "github.com/bootdotdev/bootdev/client"
)

// Whitelist commands
var allowedCommands = map[string]bool{
	"ls":   true,
	"echo": true,
	"cat":  true,
	// give more commands if needed
}

func validateArgs(args []string) bool {
	for _, arg := range args {
		for _, r := range arg {
			if !unicode.IsPrint(r) || strings.ContainsAny(string(r), `&|;`) {
				return false
			}
		}
	}
	return true
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

		parts := strings.Fields(finalCommand)
		if len(parts) == 0 {
			responses[i].ExitCode = -1
			responses[i].Stdout = "Invalid command"
			continue
		}

		// check whitelist
		if !allowedCommands[parts[0]] {
			responses[i].ExitCode = -1
			responses[i].Stdout = "Command not allowed"
			continue
		}

		if !validateArgs(parts[1:]) {
			responses[i].ExitCode = -1
			responses[i].Stdout = "Invalid arguments"
			continue
		}

		cmd := exec.Command(parts[0], parts[1:]...)

		b, err := cmd.Output()
		if ee, ok := err.(*exec.ExitError); ok {
			responses[i].ExitCode = ee.ExitCode()
		} else if err != nil {
			responses[i].ExitCode = -2
		} else {
			responses[i].ExitCode = 0
		}

		responses[i].Stdout = strings.TrimRight(string(b), " \n\t\r")
	}

	return responses
}

func interpolateArgs(rawCommand string, optionalPositionalArgs []string) string {
	for i, arg := range optionalPositionalArgs {
		rawCommand = strings.ReplaceAll(rawCommand, fmt.Sprintf("$%d", i+1), arg)
	}
	return rawCommand
}
