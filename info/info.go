package info

import (
	"fmt"
	"os"

	"github.com/launchpad-project/cli/config"
	"github.com/launchpad-project/cli/util"
	"github.com/mitchellh/mapstructure"
)

// Service structure
type Service struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var appConfig = config.Stores["app"]

func printNotEmpty(name, key string) {
	value, err := appConfig.GetString(key)

	if util.Verbose && err != nil {
		fmt.Println(key+" failed:", err)
	}

	if len(value) != 0 {
		fmt.Println(name + ": " + value)

	}
}

// GetCurrentAppInfo lists the information described on the project's launchpad.json
func GetCurrentAppInfo() {
	if len(appConfig.Data) == 0 {
		fmt.Fprintf(os.Stderr, "Application not found.\n")
		os.Exit(1)
		return
	}

	printNotEmpty("Application", "name")
	printNotEmpty("Description", "description")
	printNotEmpty("Domain", "domain")

	var i, err = appConfig.GetInterface("services")

	if err != nil {
		panic(err)
	}

	var list []Service
	mapstructure.Decode(i, &list)

	if len(list) != 0 {
		fmt.Println("\nList of services")
	}

	for _, service := range list {
		fmt.Println(service.Name, "-", service.Description)
	}

}
