package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/logrusorgru/aurora/v4"

	"github.com/spf13/cobra"
)

var flagConfirmDestroy bool

func newDestroyCommand(kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroys the current WordPress site. This is a permanent change.",
		Run: func(cmd *cobra.Command, args []string) {
			confirmDestroy := false

			if flagConfirmDestroy {
				confirmDestroy = true
			} else {
				confirmDestroy = console.PromptConfirm(
					fmt.Sprintf(
						"Are you sure you want to destroy %s? %s",
						aurora.Bold(aurora.Blue(kanaSite.Settings.Name)),
						aurora.Bold(
							aurora.Yellow(
								"This operation is destructive and cannot be undone."))),
					false)
			}

			if confirmDestroy {
				err := kanaSite.EnsureDocker()
				if err != nil {
					console.Error(err, flagVerbose)
				}

				// Stop the WordPress site.
				err = kanaSite.StopSite()
				if err != nil {
					console.Error(err, flagVerbose)
				}

				// Remove the site's folder in the config directory.
				err = os.RemoveAll(kanaSite.Settings.SiteDirectory)
				if err != nil {
					console.Error(err, flagVerbose)
				}

				console.Success(fmt.Sprintf("Your site, %s, has been completely destroyed.", aurora.Bold(aurora.Blue(kanaSite.Settings.Name))))
				return
			}

			console.Error(fmt.Errorf("site destruction canceled. No data has been lost"), flagVerbose)
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.Flags().BoolVar(&flagConfirmDestroy, "confirm-destroy", false, "Confirm destruction of your site (doesn't require a prompt).")

	return cmd
}
