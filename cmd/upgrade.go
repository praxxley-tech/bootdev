package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/bootdotdev/bootdev/version"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"update"},
	Short:   "Installs the latest version of the CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		info := version.FromContext(cmd.Context())
		if !info.IsOutdated {
			fmt.Println("Boot.dev CLI is already up to date.")
			return
		}
		// Install the latest version
		command := exec.Command("go", "install", "github.com/bootdotdev/bootdev@latest")
		output, err := command.CombinedOutput()
		if err != nil {
			fmt.Printf("Error installing latest version: %s\n", err)
			fmt.Printf("Command output: %s\n", string(output)) // Optional: Protokolliere die Ausgabe zur Fehlerbehebung
			cobra.CheckErr(err)
		}

		// Get the new version info
		command = exec.Command("bootdev", "--version")
		versionOutput, err := command.Output()
		cobra.CheckErr(err)
		re := regexp.MustCompile(`v\d+\.\d+\.\d+`)
		version := re.FindString(string(versionOutput))
		fmt.Printf("Successfully upgraded to %s!\n", version)

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
