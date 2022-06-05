package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/setup"
	"github.com/ChrisWiegman/kana/internal/traefik"
	"github.com/ChrisWiegman/kana/internal/wordpress"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kana",
	Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
	Run: func(cmd *cobra.Command, args []string) {

		kanaConfig, err := config.GetKanaConfig()
		if err != nil {
			panic(err)
		}

		controller, err := docker.NewController(kanaConfig)
		if err != nil {
			panic(err)
		}

		_, _, err = controller.EnsureNetwork("kana")
		if err != nil {
			panic(err)
		}

		setup.EnsureAppConfig(kanaConfig)
		setup.EnsureCerts(kanaConfig)
		err = traefik.NewTraefik(controller)
		if err != nil {
			panic(err)
		}

		wordpress.NewWordPress(controller)
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
