package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

func newOpenCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open the current site in your browser.",
		Run: func(cmd *cobra.Command, args []string) {
			runOpen(cmd, args, site)
		},
		Args: cobra.NoArgs,
	}

	return cmd
}

func runOpen(cmd *cobra.Command, args []string, site *site.Site) {

	// Open the site in the user's default browser,
	err := site.OpenSite()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
