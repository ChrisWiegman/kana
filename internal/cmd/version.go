package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = ""
	Timestamp = ""
)

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays version information for the Kana CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", Version)
			fmt.Printf("Build Time: %s\n", Timestamp)
		},
		Args: cobra.NoArgs,
	}

	return cmd
}
