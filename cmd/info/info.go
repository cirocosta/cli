package info

import (
	"fmt"
	"os"

	"github.com/launchpad-project/cli/launchpad/config"
	"github.com/spf13/cobra"
)

var appStore = config.Stores["app"]

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about current app directory",
	Run:   infoRun,
}

func printNotEmpty(name, key string) {
	if value := appStore.Data.GetString(key); len(value) != 0 {
		fmt.Println(name + ": " + value)
	}
}

func infoRun(cmd *cobra.Command, args []string) {
	if len(appStore.Data.GetString("name")) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", "")
		os.Exit(1)
		return
	}

	printNotEmpty("Application", "name")
	printNotEmpty("Description", "description")
	printNotEmpty("Domain", "domain")
}
