package cmddeploy

import (
	"os"

	"github.com/launchpad-project/cli/cmdcontext"
	"github.com/launchpad-project/cli/containers"
	"github.com/launchpad-project/cli/deploy"
	"github.com/spf13/cobra"
)

var output bool

// DeployCmd deploys the current project or container
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys the current project or container",
	Run:   deployRun,
}

func deployRun(cmd *cobra.Command, args []string) {
	_, _, err := cmdcontext.GetProjectOrContainerID(args)

	if err != nil {
		println("fatal: not a project")
		os.Exit(1)
	}

	list, err := containers.GetListFromScope()

	if err != nil {
		println(err)
		os.Exit(1)
	}

	switch output {
	case false:
		deploy.DeployContainers(list)
	default:
		deploy.ZipContainers(list)
	}
}

func init() {
	DeployCmd.Flags().BoolVarP(&output, "output", "o", false, "Generates a <container>.jar bundle")
}
