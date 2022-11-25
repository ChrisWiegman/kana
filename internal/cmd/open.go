package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/console"
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

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}

func runOpen(cmd *cobra.Command, args []string, site *site.Site) {

	// Open the site in the user's default browser,
	err := site.OpenSite()
	if err != nil {
		console.Error(fmt.Errorf("the site doesn't appear to be running. Please use `kana start` to start the site"), flagDebugMode)
	}
}
