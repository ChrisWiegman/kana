package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/spf13/cobra"
)

var openAppFlag string

func newOpenCommand(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open the current site in your browser.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			// Open the site in the user's default browser,
			err = kanaSite.OpenSite(openAppFlag)
			if err != nil {
				consoleOutput.Error(fmt.Errorf("the site doesn't appear to be running. Please use `kana start` to start the site"))
			}

			consoleOutput.Success(
				fmt.Sprintf(
					"Your site, %s, has been opened in your default browser.",
					consoleOutput.Bold(
						consoleOutput.Blue(
							kanaSite.Settings.Name))))
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.Flags().StringVarP(&openAppFlag, "app", "a", "site", "site = open kana site, phpmyadmin = PhpMyAdmin, mailpit = Mailpit")

	return cmd
}
