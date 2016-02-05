package auth

import (
	"fmt"

	"github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/launchpad"
	"github.com/launchpad-project/cli/prompt"
	"github.com/spf13/cobra"
)

var globalConfig = config.Stores["global"]

// LoginCmd sets the user credential
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Using Basic Authentication with your credentials",
	Run:   loginRun,
}

// LogoutCmd unsets the user credential
var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Revoke credentials",
	Run:   logoutRun,
}

func loginRun(cmd *cobra.Command, args []string) {
	var username = prompt.Prompt("Username")
	var password = prompt.Prompt("Password")

	globalConfig.Set("endpoint", launchpad.Endpoint)
	globalConfig.Set("username", username)
	globalConfig.Set("password", password)
	globalConfig.Save()

	fmt.Println("Authentication information saved.")
}

func logoutRun(cmd *cobra.Command, args []string) {
	globalConfig.Set("username", "")
	globalConfig.Set("password", "")
	globalConfig.Save()
}
