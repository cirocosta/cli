package config

import (
	"fmt"
	"os"
	"strings"

	c "github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/prompt"
	"github.com/spf13/cobra"
)

var (
	listParam   bool
	setParam    bool
	globalParam bool
	configStore = c.Stores["app"]
)

// ConfigCmd is used for configuring the CLI tool and app
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration for the Launchpad CLI tool",
	Run:   configRun,
}

func listKeys() {
	for key, configurable := range configStore.ConfigurableKeys {
		if configurable {
			fmt.Println(key, "=", configStore.Get(key))
		}
	}
}

func configRun(cmd *cobra.Command, args []string) {
	if globalParam {
		configStore = c.Stores["global"]
	}

	if listParam {
		listKeys()
		return
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
		configStore.SetAndSavePublicKey(key, prompt.Prompt(key))
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
	ConfigCmd.Flags().BoolVarP(&listParam, "list", "l", false, "list all")
	ConfigCmd.Flags().BoolVar(&setParam, "set", false, "set property")
	ConfigCmd.Flags().BoolVar(&globalParam, "global", false, "set application global config property")
}
