package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/logrusorgru/aurora/v4"

	"github.com/spf13/cobra"
)

func newOpenCommand(kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open the current site in your browser.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// Open the site in the user's default browser,
			err = kanaSite.OpenSite()
			if err != nil {
				console.Error(fmt.Errorf("the site doesn't appear to be running. Please use `kana start` to start the site"), flagVerbose)
			}

			console.Success(fmt.Sprintf("Your site, %s, has been opened in your default browser.", aurora.Bold(aurora.Blue(kanaSite.Settings.Name))))
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}
