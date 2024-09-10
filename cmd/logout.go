package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func logout() {
	api_url := viper.GetString("api_url")
	client := &http.Client{}

	// Best effort - logout should never fail, but handle errors if they occur
	r, err := http.NewRequest("POST", api_url+"/v1/auth/logout", bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fehler beim Erstellen der Logout-Anfrage: %v\n", err)
		return
	}
	r.Header.Add("X-Refresh-Token", viper.GetString("refresh_token"))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fehler beim Senden der Logout-Anfrage: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Logout-Anfrage fehlgeschlagen mit Statuscode: %d\n", resp.StatusCode)
		return
	}

	viper.Set("access_token", "")
	viper.Set("refresh_token", "")
	viper.Set("last_refresh", time.Now().Unix())
	if err := viper.WriteConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Fehler beim Schreiben der Konfigurationsdatei: %v\n", err)
		return
	}

	fmt.Println("Logged out successfully.")
}

var logoutCmd = &cobra.Command{
	Use:          "logout",
	Aliases:      []string{"signout"},
	Short:        "Disconnect the CLI from your account",
	PreRun:       requireAuth,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logout()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
