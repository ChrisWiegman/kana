package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

var (
	Version   = ""
	GitHash   = ""
	Timestamp = ""
)

func newVersionCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays version information for the Kana CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			runVersion(cmd, args, site)
		},
		Args: cobra.NoArgs,
	}

	return cmd
}

func runVersion(cmd *cobra.Command, args []string, site *site.Site) {

	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Commit Hash: %s\n", GitHash)
	fmt.Printf("Build Time: %s\n", Timestamp)
}
