package cmd

import (
	"os"

	"github.com/launchpad-project/cli/cmd/auth"
	cconfig "github.com/launchpad-project/cli/cmd/config"
	"github.com/launchpad-project/cli/cmd/info"
	"github.com/launchpad-project/cli/cmd/services"
	cupdate "github.com/launchpad-project/cli/cmd/update"
	"github.com/launchpad-project/cli/cmd/version"
	"github.com/launchpad-project/cli/launchpad"
	"github.com/launchpad-project/cli/launchpad/config"
	"github.com/launchpad-project/cli/launchpad/update"
	"github.com/launchpad-project/cli/launchpad/util"
	"github.com/spf13/cobra"
)

var Verbose bool

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

func preRun(cmd *cobra.Command, args []string) {
	util.Verbose = Verbose

	update.PostUpdate()

	// if != launchpad login or launchpad config and not logged, exit with error
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	config.Setup()
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	RootCmd.AddCommand(info.InfoCmd)
	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(services.ServicesCmd)
	RootCmd.AddCommand(cupdate.UpdateCmd)
	RootCmd.AddCommand(version.VersionCmd)
	RootCmd.AddCommand(cconfig.ConfigCmd)
}
