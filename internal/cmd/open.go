package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/spf13/cobra"
)

var openPhpMyAdminFlag, openMailpitFlag, openSiteFlag bool

func newOpenCommand(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open the current site in your browser.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			// Default to opening the site if no flags are specified
			if !cmd.Flags().Lookup("phpmyadmin").Changed &&
				!cmd.Flags().Lookup("mailpit").Changed &&
				!cmd.Flags().Lookup("site").Changed {
				openSiteFlag = true
			}

			// Open the site in the user's default browser,
			err = kanaSite.OpenSite(openPhpMyAdminFlag, openMailpitFlag, openSiteFlag, consoleOutput)
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
	cmd.Flags().BoolVarP(
		&openPhpMyAdminFlag,
		"phpmyadmin",
		"p",
		false,
		"Opens the PhpMyAdmin UI for the current or specified Kana site in your default browser")
	cmd.Flags().BoolVarP(
		&openMailpitFlag,
		"mailpit",
		"m",
		false,
		"Opens the Mailpit UI for the current or specified Kana site in your default browser")
	cmd.Flags().BoolVarP(&openSiteFlag, "site", "s", false, "Opens the current or specified Kana site in your default browser")

	return cmd
}
