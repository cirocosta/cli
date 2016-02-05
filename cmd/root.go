package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/launchpad-project/cli/cmd/auth"
	cconfig "github.com/launchpad-project/cli/cmd/config"
	"github.com/launchpad-project/cli/cmd/info"
	"github.com/launchpad-project/cli/cmd/services"
	cupdate "github.com/launchpad-project/cli/cmd/update"
	"github.com/launchpad-project/cli/cmd/version"
	"github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/launchpad"
	"github.com/launchpad-project/cli/util"
	"github.com/spf13/cobra"
)

var noAuthWhitelist = map[string]bool{
	"login":  true,
	"logout": true,
	"config": true,
	"update": true,
}

// RootCmd is the main command for the CLI
var RootCmd = &cobra.Command{
	Use:   "launchpad",
	Short: "Launchpad CLI tool",
	Long: `Launchpad Command Line Interface
Version ` + launchpad.Version + `
Copyright 2016 Liferay, Inc.
http://liferay.io`,
	PersistentPreRun: preRun,
}

var globalStore = config.Stores["global"]

func verifyAuth(commandPath string) {
	var test = strings.SplitAfterN(commandPath, " ", 2)[1]

	for key := range noAuthWhitelist {
		if key == test {
			return
		}
	}

	_, err1 := globalStore.GetString("endpoint")
	_, err2 := globalStore.GetString("username")
	_, err3 := globalStore.GetString("password")

	if err1 == nil && err2 == nil && err3 == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Please run \"launchpad login\" first.\n")
	os.Exit(1)
}

func preRun(cmd *cobra.Command, args []string) {
	verifyAuth(cmd.CommandPath())
}

// Execute is the Entry-point for the CLI
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	config.Setup()
	RootCmd.PersistentFlags().BoolVarP(&util.Verbose, "verbose", "v", false, "verbose output")
	RootCmd.AddCommand(info.InfoCmd)
	RootCmd.AddCommand(auth.LoginCmd)
	RootCmd.AddCommand(auth.LogoutCmd)
	RootCmd.AddCommand(services.ServicesCmd)
	RootCmd.AddCommand(cupdate.UpdateCmd)
	RootCmd.AddCommand(version.VersionCmd)
	RootCmd.AddCommand(cconfig.ConfigCmd)
}
