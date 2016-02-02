package auth

import (
	"fmt"

	"github.com/launchpad-project/cli/launchpad"
	"github.com/launchpad-project/cli/launchpad/config"
	"github.com/spf13/cobra"
)

var store = config.Stores["app"]

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication for the Launchpad CLI tool",
}

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

	store.Data.Set("endpoint", "https://liferay.io/")
	store.Data.Set("username", username)
	store.Data.Set("password", password)
	store.Save()

	fmt.Println("Authentication information saved.")
}

func logoutRun(cmd *cobra.Command, args []string) {
	store.Data.Set("username", "")
	store.Data.Set("password", "")
	store.Save()
}

func init() {
	AuthCmd.AddCommand(LoginCmd)
	AuthCmd.AddCommand(LogoutCmd)
}
