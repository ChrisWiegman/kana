package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/docker"

	"github.com/spf13/cobra"
)

func newRootCommand() (*cobra.Command, *docker.Controller) {

	kanaConfig, err := config.GetKanaConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	controller, err := docker.NewController(kanaConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   "kana",
		Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("")
		},
	}

	return cmd, controller

}

func Execute() {

	cmd, controller := newRootCommand()

	cmd.AddCommand(
		newStartCommand(controller),
		newStopCommand(controller),
		newOpenCommand(controller),
	)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
