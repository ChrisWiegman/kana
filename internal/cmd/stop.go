package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/logrusorgru/aurora/v4"

	"github.com/spf13/cobra"
)

func newStopCommand(kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops the WordPress development environment.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// Stop the WordPress site
			err = kanaSite.StopSite()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			console.Success(
				fmt.Sprintf(
					"Your site, %s, has been stopped. Please run `kana start` again if you would like to use it.",
					aurora.Bold(aurora.Blue(kanaSite.Settings.Name))))
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}
