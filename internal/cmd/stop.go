package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/config"
	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/logrusorgru/aurora/v4"

	"github.com/spf13/cobra"
)

func newStopCommand(kanaConfig *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops the WordPress development environment.",
		Run: func(cmd *cobra.Command, args []string) {

			site, err := site.NewSite(kanaConfig)
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// Stop the WordPress site
			err = site.StopWordPress()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			console.Success(fmt.Sprintf("Your site, %s, has been stopped. Please run `kana start` again if you would like to use it.", aurora.Bold(aurora.Blue(kanaConfig.Site.Name))))
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}
