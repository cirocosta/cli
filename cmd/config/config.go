package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/launchpad-project/cli/launchpad"
	c "github.com/launchpad-project/cli/launchpad/config"
	"github.com/spf13/cobra"
)

var (
	setParam    bool
	globalParam bool
	configStore = c.Stores["app"]
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration for the Launchpad CLI tool",
	Run:   configRun,
}

func configRun(cmd *cobra.Command, args []string) {
	if globalParam {
		configStore = c.Stores["global"]
	}

	if len(args) == 0 {
		cmd.Help()
		return
	}

	var key = args[0]

	if len(args) != 1 {
		configStore.SetAndSavePublicKey(key, strings.Join(args[1:], " "))
		return
	}

	if setParam {
		configStore.SetAndSavePublicKey(key, launchpad.Prompt(key))
		return
	}

	var value, err = configStore.GetString(key)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Println(value)
}

func init() {
	ConfigCmd.Flags().BoolVar(&setParam, "set", false, "Set property")
	ConfigCmd.Flags().BoolVar(&globalParam, "global", false, "Set application global config property")
}
