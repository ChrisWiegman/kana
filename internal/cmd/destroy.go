package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var flagForce bool

func newDestroyCommand(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroys the current WordPress site. This is a permanent change.",
		Run: func(cmd *cobra.Command, args []string) {
			confirmDestroy := false

			if flagForce {
				confirmDestroy = true
			} else {
				confirmDestroy = consoleOutput.PromptConfirm(
					fmt.Sprintf(
						"Are you sure you want to destroy %s? %s",
						consoleOutput.Bold(consoleOutput.Blue(kanaSite.Settings.Name)),
						consoleOutput.Bold(
							consoleOutput.Yellow(
								"This operation is destructive and cannot be undone."))),
					false)
			}

			if confirmDestroy {
				err := kanaSite.EnsureDocker(consoleOutput)
				if err != nil {
					consoleOutput.Error(err)
				}

				// Stop the WordPress site.
				err = kanaSite.StopSite()
				if err != nil {
					consoleOutput.Error(err)
				}

				// Remove the site's folder in the config directory.
				err = os.RemoveAll(kanaSite.Settings.SiteDirectory)
				if err != nil {
					consoleOutput.Error(err)
				}

				consoleOutput.Success(
					fmt.Sprintf(
						"Your site, %s, has been completely destroyed.",
						consoleOutput.Bold(
							consoleOutput.Blue(
								kanaSite.Settings.Name))))
				return
			}

			consoleOutput.Error(fmt.Errorf("site destruction canceled. No data has been lost"))
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.Flags().BoolVar(&flagForce, "force", false, "Force destruction of your site (doesn't require a prompt).")
	cmd.Flags().SetNormalizeFunc(aliasForceFlag)
	return cmd
}

func aliasForceFlag(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if name == "confirm-destroy" {
		name = "force"
	}

	return pflag.NormalizedName(name)
}
