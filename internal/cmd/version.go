package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	Version   = ""
	GitHash   = ""
	Timestamp = ""
)

func newVersionCommand(kanaConfig *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays version information for the Kana CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			runVersion(cmd, args, kanaConfig)
		},
		Args: cobra.NoArgs,
	}

	return cmd
}

func runVersion(cmd *cobra.Command, args []string, kanaConfig *config.Config) {

	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Commit Hash: %s\n", GitHash)
	fmt.Printf("Build Time: %s\n", Timestamp)
}
