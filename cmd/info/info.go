package info

import (
	"github.com/launchpad-project/cli/launchpad/config"
	"github.com/launchpad-project/cli/launchpad/info"
	"github.com/spf13/cobra"
)

var appConfig = config.Stores["app"]

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about current app directory",
	Run:   infoRun,
}

func infoRun(cmd *cobra.Command, args []string) {
	info.GetCurrentAppInfo()
}
