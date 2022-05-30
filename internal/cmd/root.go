package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/setup"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kana",
	Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
	Run: func(cmd *cobra.Command, args []string) {

		setup.EnsureCA()
		setup.EnsureAppConfig()
		//traefik.NewTraefik()
		//wordpress.NewWordPress()
	},
}

func Execute() {

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
