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

		controller.EnsureImage("wordpress")
		statusCode, err := controller.ContainerRun("wordpress", []string{}, []docker.VolumeMount{})
		hasNetwork, network, _ := controller.EnsureNetwork("kana")

		if hasNetwork {
			fmt.Println(network.Name)
		} else {
			fmt.Println("Network doesn't exist")
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Println(statusCode)

	},
}

func Execute() {

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
