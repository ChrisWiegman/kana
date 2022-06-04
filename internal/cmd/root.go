package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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

		controller, err := docker.NewController()
		if err != nil {
			panic(err)
		}

		_, _, err = controller.EnsureNetwork("kana")
		if err != nil {
			panic(err)
		}

		setup.EnsureAppConfig()
		setup.EnsureCerts()
		err = traefik.NewTraefik(controller)
		if err != nil {
			panic(err)
		}

		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		wordpress.NewWordPress(filepath.Base(path), controller)
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
