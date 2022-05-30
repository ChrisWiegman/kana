package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/pkg/minica"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kana",
	Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
	Run: func(cmd *cobra.Command, args []string) {

		minica.GenCerts("/Users/chriswiegman/.kana/certs/kana.key", "/Users/chriswiegman/.kana/certs/kana.pem", []string{"wordpress.local.dev"}, []string{})
	},
}

func Execute() {

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
