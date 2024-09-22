package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-wordpress/internal/console"
	"github.com/ChrisWiegman/kana-wordpress/internal/settings"
	"github.com/ChrisWiegman/kana-wordpress/internal/site"

	"github.com/spf13/cobra"
)

func stop(consoleOutput *console.Console, kanaSite *site.Site, kanaSettings *settings.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops the WordPress development environment.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			// Stop the WordPress site
			err = kanaSite.StopSite()
			if err != nil {
				consoleOutput.Error(err)
			}

			consoleOutput.Success(
				fmt.Sprintf(
					"Your site, %s, has been stopped. Please use `kana start` again to restart it.",
					consoleOutput.Bold(consoleOutput.Blue(kanaSettings.Get("name")))))
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}
