package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wedeploy/cli/apihelper"
	"github.com/wedeploy/cli/cmd/auth"
	"github.com/wedeploy/cli/cmd/build"
	"github.com/wedeploy/cli/cmd/createctx"
	"github.com/wedeploy/cli/cmd/link"
	"github.com/wedeploy/cli/cmd/list"
	"github.com/wedeploy/cli/cmd/logs"
	"github.com/wedeploy/cli/cmd/remote"
	"github.com/wedeploy/cli/cmd/restart"
	"github.com/wedeploy/cli/cmd/run"
	"github.com/wedeploy/cli/cmd/stop"
	"github.com/wedeploy/cli/cmd/unlink"
	"github.com/wedeploy/cli/cmd/update"
	"github.com/wedeploy/cli/cmd/version"
	"github.com/wedeploy/cli/color"
	"github.com/wedeploy/cli/config"
	"github.com/wedeploy/cli/defaults"
	"github.com/wedeploy/cli/verbose"
	"github.com/wedeploy/cli/verbosereq"
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

// ListNoRemoteFlags hides the globals non used --remote
var ListNoRemoteFlags = map[string]bool{
	"link":    true,
	"unlink":  true,
	"run":     true,
	"stop":    true,
	"remote":  true,
	"update":  true,
	"version": true,
}

// LocalOnlyCommands for local-only commands
var LocalOnlyCommands = map[string]bool{
	"link":   true,
	"unlink": true,
	"run":    true,
	"stop":   true,
}

// RootCmd is the main command for the CLI
var RootCmd = &cobra.Command{
	Use:   "we",
	Short: "WeDeploy CLI tool",
	Long: `WeDeploy Command Line Interface
Version ` + defaults.Version + `
Copyright 2016 Liferay, Inc.
http://wedeploy.com`,
	PersistentPreRunE: persistentPreRun,
	Run:               run,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var (
	version bool
	remote  string
)

var commands = []*cobra.Command{
	cmdauth.LoginCmd,
	cmdauth.LogoutCmd,
	cmdcreate.CreateCmd,
	cmdlogs.LogsCmd,
	cmdlist.ListCmd,
	cmdrestart.RestartCmd,
	cmdbuild.BuildCmd,
	cmdrun.RunCmd,
	cmdstop.StopCmd,
	cmdlink.LinkCmd,
	cmdunlink.UnlinkCmd,
	cmdremote.RemoteCmd,
	cmdupdate.UpdateCmd,
	cmdversion.VersionCmd,
}

func hideVersionFlag() {
	if err := RootCmd.Flags().MarkHidden("version"); err != nil {
		panic(err)
	}
}

func hideNoVerboseRequestsFlag() {
	if err := RootCmd.PersistentFlags().MarkHidden("no-verbose-requests"); err != nil {
		panic(err)
	}
}

func hideUnusedGlobalRemoteFlags() {
	var args = os.Args

	if len(args) < 2 {
		return
	}

	_, h := ListNoRemoteFlags[args[1]]

	if !h {
		return
	}

	if err := RootCmd.PersistentFlags().MarkHidden("remote"); err != nil {
		panic(err)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(
		&verbose.Enabled,
		"verbose",
		"v",
		false,
		"Verbose output")

	RootCmd.PersistentFlags().BoolVar(
		&verbosereq.Disabled,
		"no-verbose-requests",
		false,
		"Hide verbose output for requests")

	RootCmd.PersistentFlags().BoolVar(
		&color.NoColor,
		"no-color",
		false,
		"Disable color output")

	RootCmd.PersistentFlags().StringVar(
		&remote,
		"remote", "", "Remote to use")

	RootCmd.Flags().BoolVar(
		&version,
		"version", false, "Print version information and quit")

	hideVersionFlag()
	hideNoVerboseRequestsFlag()
	hideUnusedGlobalRemoteFlags()

	for _, c := range commands {
		RootCmd.AddCommand(c)
	}
}

func setEndpoint() error {
	if isLocalCommandOnly() && remote != "" {
		return errors.New("can not use command with a remote")
	}

	if remote == "" {
		return setLocal()
	}

	return setRemote()
}

func setLocal() error {
	config.Context.Token = apihelper.DefaultToken
	config.Context.Endpoint = fmt.Sprintf("http://localhost:%d/", config.Global.LocalPort)
	return nil
}

func setRemote() error {
	var r, ok = config.Global.Remotes.Get(remote)

	if !ok {
		return errors.New("Remote " + remote + " is not configured.")
	}

	config.Context.Remote = remote
	config.Context.Endpoint = normalizeRemote(r.URL)
	config.Context.Username = config.Global.Username
	config.Context.Password = config.Global.Password
	config.Context.Token = config.Global.Token
	return nil
}

func normalizeRemote(address string) string {
	if address != "" &&
		!strings.HasPrefix(address, "http://") &&
		!strings.HasPrefix(address, "https://") {
		return "http://" + address
	}

	return address
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	if config.Global.NoColor {
		color.NoColor = true
	}

	if err := verifyCmdReqAuth(cmd.CommandPath()); err != nil {
		return err
	}

	return setEndpoint()
}

func run(cmd *cobra.Command, args []string) {
	if version {
		cmdversion.VersionCmd.Run(cmd, args)
		return
	}

	if err := cmd.Help(); err != nil {
		panic(err)
	}
}

func isCmdWhitelistNoAuth(commandPath string) bool {
	var parts = strings.SplitAfterN(commandPath, " ", 2)

	if len(parts) < 2 {
		return true
	}

	for key := range WhitelistCmdsNoAuthentication {
		if key == parts[1] {
			return true
		}
	}

	return false
}

func isLocalCommandOnly() bool {
	var args = os.Args

	if len(args) < 2 {
		return false
	}

	_, h := LocalOnlyCommands[args[1]]
	return h
}

func verifyCmdReqAuth(commandPath string) error {
	if remote == "" || isCmdWhitelistNoAuth(commandPath) {
		return nil
	}

	var g = config.Global

	var hasAuth = (g.Token != "") || (g.Username != "" && g.Password != "")

	if g.Endpoint != "" && hasAuth {
		return nil
	}

	return errors.New(`Please run "we login" first.`)
}
