package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/logrusorgru/aurora/v4"

	"github.com/spf13/cobra"
)

func newStopCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops the WordPress development environment.",
		Run: func(cmd *cobra.Command, args []string) {
			runStop(cmd, args, site)
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}

func runStop(cmd *cobra.Command, args []string, site *site.Site) {

	// Stop the WordPress site
	err := site.StopWordPress()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	console.Success(fmt.Sprintf("Your site, %s, has been stopped. Please run `kana start` again if you would like to use it.", aurora.Bold(aurora.Blue(site.StaticConfig.SiteName))))
}
