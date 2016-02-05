package version

import (
	"fmt"
	"runtime"

	"github.com/launchpad-project/cli/launchpad"
	"github.com/spf13/cobra"
)

// VersionCmd is used for reading the version of this tool
var VersionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		var os = runtime.GOOS
		var arch = runtime.GOARCH
		fmt.Printf("Launchpad CLI version %s %s/%s\n", launchpad.Version, os, arch)
	},
	Short: "Prints the Command Line Tools version",
}
