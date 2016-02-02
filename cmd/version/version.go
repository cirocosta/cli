package version

import (
	"fmt"

	"github.com/launchpad-project/cli/launchpad"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Launchpad CLI version", launchpad.Version)
	},
	Short: "Prints the Command Line Tools version",
}
