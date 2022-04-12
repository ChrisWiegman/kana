package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kana",
	Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This is Kana.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
