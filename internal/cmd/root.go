package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/m/pkg/docker"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kana",
	Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
	Run: func(cmd *cobra.Command, args []string) {
		var controller, _ = docker.NewController()

		controller.EnsureImage("alpine")
		statusCode, body, err := controller.ContainerRunAndClean("alpine", []string{"echo", "hello world"}, []docker.VolumeMount{})

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Println(statusCode)
		fmt.Println(body)

	},
}

func Execute() {

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
