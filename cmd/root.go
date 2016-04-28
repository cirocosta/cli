package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/launchpad-project/cli/cmd/auth"
	cmdcontainers "github.com/launchpad-project/cli/cmd/containers"
	cmdcreate "github.com/launchpad-project/cli/cmd/createctx"
	cmddeploy "github.com/launchpad-project/cli/cmd/deploy"
	cmdhooks "github.com/launchpad-project/cli/cmd/hooks"
	cmdlogs "github.com/launchpad-project/cli/cmd/logs"
	cmdprojects "github.com/launchpad-project/cli/cmd/projects"
	cmdrestart "github.com/launchpad-project/cli/cmd/restart"
	cmdstatus "github.com/launchpad-project/cli/cmd/status"
	cmdupdate "github.com/launchpad-project/cli/cmd/update"
	cmdversion "github.com/launchpad-project/cli/cmd/version"
	"github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/configstore"
	"github.com/launchpad-project/cli/defaults"
	"github.com/launchpad-project/cli/run"
	"github.com/launchpad-project/cli/update"
	"github.com/launchpad-project/cli/verbose"
	"github.com/spf13/cobra"
)

// WhitelistCmdsNoAuthentication for cmds that doesn't require authentication
var WhitelistCmdsNoAuthentication = map[string]bool{
	"login":   true,
	"logout":  true,
	"build":   true,
	"deploy":  true,
	"update":  true,
	"version": true,
}

// RootCmd is the main command for the CLI
var RootCmd = &cobra.Command{
	Use:   "launchpad",
	Short: "Launchpad CLI tool",
	Long: `Launchpad Command Line Interface
Version ` + defaults.Version + `
Copyright 2016 Liferay, Inc.
http://liferay.io`,
	PersistentPreRun: persistentPreRun,
}

var globalStore *configstore.Store

// Execute is the Entry-point for the CLI
func Execute() {
	var wgUpdate sync.WaitGroup
	var errUpdate error

	wgUpdate.Add(1)
	go func() {
		errUpdate = update.NotifierCheck()
		wgUpdate.Done()
	}()

	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}

	wgUpdate.Wait()

	if errUpdate == nil {
		update.Notify()
	} else {
		println("Update notification error:", errUpdate.Error())
	}
}

func init() {
	config.Setup()

	RootCmd.PersistentFlags().BoolVarP(
		&verbose.Enabled,
		"verbose",
		"v",
		false,
		"verbose output")

	RootCmd.PersistentFlags().BoolVar(
		&color.NoColor,
		"no-color",
		false,
		"disable color output")

	var csg = config.Stores["global"]

	if csg.Get("no_color") == "true" {
		color.NoColor = true
	}

	RootCmd.AddCommand(cmdauth.LoginCmd)
	RootCmd.AddCommand(cmdauth.LogoutCmd)
	RootCmd.AddCommand(cmdcreate.CreateCmd)
	RootCmd.AddCommand(cmdlogs.LogsCmd)
	RootCmd.AddCommand(cmdprojects.ProjectsCmd)
	RootCmd.AddCommand(cmdcontainers.ContainersCmd)
	RootCmd.AddCommand(cmdstatus.StatusCmd)
	RootCmd.AddCommand(cmdrestart.RestartCmd)
	RootCmd.AddCommand(cmdhooks.BuildCmd)
	RootCmd.AddCommand(cmdrun.RunCmd)
	RootCmd.AddCommand(cmddeploy.DeployCmd)
	RootCmd.AddCommand(cmdupdate.UpdateCmd)
	RootCmd.AddCommand(cmdversion.VersionCmd)
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	verifyAuth(cmd.CommandPath())
}

func verifyAuth(commandPath string) {
	var csg = config.Stores["global"]
	var test = strings.SplitAfterN(commandPath, " ", 2)[1]

	for key := range WhitelistCmdsNoAuthentication {
		if key == test {
			return
		}
	}

	_, err1 := csg.GetString("endpoint")
	_, err2 := csg.GetString("username")
	_, err3 := csg.GetString("password")

	if err1 == nil && err2 == nil && err3 == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Please run \"launchpad login\" first.\n")
	os.Exit(1)
}
