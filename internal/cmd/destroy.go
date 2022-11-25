package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/logrusorgru/aurora/v4"

	"github.com/spf13/cobra"
)

func newDestroyCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroys the current WordPress site. This is a permanent change.",
		Run: func(cmd *cobra.Command, args []string) {
			runDestroy(cmd, args, site)
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}

func runDestroy(cmd *cobra.Command, args []string, site *site.Site) {

	// Stop the WordPress site.
	err := site.StopWordPress()
	if err != nil {
		console.Error(err, flagDebugMode)
	}

	// Remove the site's folder in the config directory.
	err = os.RemoveAll(site.StaticConfig.SiteDirectory)
	if err != nil {
		console.Error(err, flagDebugMode)
	}

	console.Success(fmt.Sprintf("Your site, %s, has been completely destroyed.", aurora.Bold(aurora.Blue(site.StaticConfig.SiteName))))
}
