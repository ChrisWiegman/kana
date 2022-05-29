package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/pkg/docker"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kana",
	Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
	Run: func(cmd *cobra.Command, args []string) {
		controller, err := docker.NewController()
		if err != nil {
			panic(err)
		}

		_, _, err = controller.EnsureNetwork("kana")
		if err != nil {
			panic(err)
		}

		err = controller.EnsureImage("caddy")
		if err != nil {
			panic(err)
		}

		_, err = controller.ContainerRun("caddy", []string{}, []docker.VolumeMount{})
		if err != nil {
			panic(err)
		}

	},
}

func Execute() {

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
