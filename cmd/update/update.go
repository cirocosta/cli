package update

import (
	"fmt"

	"github.com/launchpad-project/cli/launchpad"
	"github.com/launchpad-project/cli/launchpad/update"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Run:   UpdateRun,
	Short: "Updates this tool to the latest version",
}

func UpdateRun(cmd *cobra.Command, args []string) {
	fmt.Println("Trying to update Launchpad CLI")
	fmt.Println("Current installed version is " + launchpad.Version)
	update.UpdateToLatest()
}
