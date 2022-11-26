package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/config"
	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/logrusorgru/aurora/v4"

	"github.com/spf13/cobra"
)

func newOpenCommand(site *site.Site, kanaConfig *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open the current site in your browser.",
		Run: func(cmd *cobra.Command, args []string) {
			runOpen(cmd, args, site, kanaConfig)
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}

func runOpen(cmd *cobra.Command, args []string, site *site.Site, kanaConfig *config.Config) {

	fmt.Println(kanaConfig.App.PHP)
	fmt.Println(site.DynamicConfig.GetString("php"))
	fmt.Println(kanaConfig.Directories.App)
	fmt.Println(site.StaticConfig.AppDirectory)
	fmt.Println(kanaConfig.Site.PHP)
	fmt.Println(site.SiteConfig.GetString("php"))
	fmt.Println(kanaConfig.Site.SiteName)
	fmt.Println(site.StaticConfig.SiteName)
	fmt.Println(kanaConfig.Directories.Site)
	fmt.Println(site.StaticConfig.SiteDirectory)
	fmt.Println(kanaConfig.Directories.Working)
	fmt.Println(site.StaticConfig.WorkingDirectory)

	return

	// Open the site in the user's default browser,
	err := site.OpenSite()
	if err != nil {
		console.Error(fmt.Errorf("the site doesn't appear to be running. Please use `kana start` to start the site"), flagVerbose)
	}

	console.Success(fmt.Sprintf("Your site, %s, has been opened in your default browser.", aurora.Bold(aurora.Blue(site.StaticConfig.SiteName))))
}
