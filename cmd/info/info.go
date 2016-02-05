package info

import (
	"github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/info"
	"github.com/spf13/cobra"
)

var appConfig = config.Stores["app"]

// InfoCmd is used for getting info about an app
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about current app directory",
	Run:   infoRun,
}

func infoRun(cmd *cobra.Command, args []string) {
	info.GetCurrentAppInfo()
}
