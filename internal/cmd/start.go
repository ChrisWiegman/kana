package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/spf13/cobra"
)

func newStartCommand(controller *docker.Controller) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new environment in the local folder.",
		Run: func(cmd *cobra.Command, args []string) {
			runStart(cmd, args, controller)
		},
	}

	return cmd

}

func runStart(cmd *cobra.Command, args []string, controller *docker.Controller) {

	fmt.Println("This is where we start the container.")

}
