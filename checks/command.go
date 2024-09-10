package checks

import (
	"fmt"
	"os/exec"
	"strings"

	api "github.com/bootdotdev/bootdev/client"
)

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
	// replace $1, $2, etc. with the optional positional args
	for i, arg := range optionalPositionalArgs {
		rawCommand = strings.ReplaceAll(rawCommand, fmt.Sprintf("$%d", i+1), arg)
	}
	return rawCommand
}
