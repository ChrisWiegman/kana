package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/wordpress"

	"github.com/spf13/cobra"
)

func newStopCommand(controller *docker.Controller) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops the WordPress development environment.",
		Run: func(cmd *cobra.Command, args []string) {
			runStop(cmd, args, controller)
		},
	}

	return cmd

}

func runStop(cmd *cobra.Command, args []string, controller *docker.Controller) {

	err := wordpress.StopWordPress(controller)
	if err != nil {
		fmt.Println(err)
	}

}
