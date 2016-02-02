package services

import (
	"github.com/launchpad-project/cli/launchpad/services"
	"github.com/spf13/cobra"
)

var ServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Manage Launchpad services",
}

var PodServicesCmd = &cobra.Command{
	Use:   "pod",
	Short: "List pods",
	Run:   PodServicesRun,
}

func PodServicesRun(cmd *cobra.Command, args []string) {
	services.GetPods()
}

func init() {
	ServicesCmd.AddCommand(PodServicesCmd)
}
