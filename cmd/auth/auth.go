package auth

import (
	"fmt"

	"github.com/launchpad-project/cli/launchpad"
	"github.com/launchpad-project/cli/launchpad/config"
	"github.com/spf13/cobra"
)

var globalConfig = config.Stores["global"]

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Using Basic Authentication with your credentials",
	Run:   loginRun,
}

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Revoke credentials",
	Run:   logoutRun,
}

func loginRun(cmd *cobra.Command, args []string) {
	var username = launchpad.Prompt("Username")
	var password = launchpad.Prompt("Password")

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
