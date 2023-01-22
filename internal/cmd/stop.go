package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/spf13/cobra"
)

func newStopCommand(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
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
					"Your site, %s, has been stopped. Please run `kana start` again if you would like to use it.",
					consoleOutput.Bold(consoleOutput.Blue(kanaSite.Settings.Name))))
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}
