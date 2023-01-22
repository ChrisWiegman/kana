package cmd

import (
	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/spf13/cobra"
)

var Version, Timestamp string

func newVersionCommand(consoleOutput *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays version information for the Kana CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			consoleOutput.Printf("Version: %s\n", Version)
			consoleOutput.Printf("Build Time: %s\n", Timestamp)
		},
		Args: cobra.NoArgs,
	}

	return cmd
}
