package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/settings"
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
)

var startFlags settings.StartFlags

func newStartCommand(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new environment in the local folder.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			fmt.Println(startFlags)
			fmt.Println(startFlags.Multisite)
			os.Exit(0)

			// Check that the site is already running and show an error if it is.
			if kanaSite.IsSiteRunning() {
				consoleOutput.Error(fmt.Errorf("the site is already running. Please stop your site before running the start command"))
			}

			// Check that we're not using our home directory as the working directory as that could cause security or other issues.
			home, err := homedir.Dir()
			if err != nil {
				consoleOutput.Error(err)
			}

			if home == kanaSite.Settings.WorkingDirectory {
				// Remove the site's folder in the config directory.
				err = os.RemoveAll(kanaSite.Settings.SiteDirectory)
				if err != nil {
					consoleOutput.Error(err)
				}

				consoleOutput.Error(fmt.Errorf("you are attempting to start a new site from your home directory. This could create security issues. Please create a folder and start a site from there")) //nolint:lll
			}

			err = kanaSite.StartSite(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			consoleOutput.Success(
				fmt.Sprintf(
					"Your site, %s, has has started and should be open in your default browser.",
					consoleOutput.Bold(consoleOutput.Blue(kanaSite.Settings.Name))))
		},
		Args: cobra.NoArgs,
	}

	// Add associated flags to customize the site at runtime.
	cmd.Flags().BoolVarP(&startFlags.Xdebug, "xdebug", "x", false, "Enable Xdebug when starting the container.")
	cmd.Flags().BoolVarP(&startFlags.WPDebug, "wpdebug", "d", false, "Enable WP_Debug when starting the container.")
	cmd.Flags().BoolVarP(&startFlags.Mailpit, "mailpit", "m", false, "Enable Mailpit when starting the container.")
	cmd.Flags().BoolVarP(&startFlags.IsPlugin, "plugin", "p", false, "Run the site as a plugin using the current folder as the plugin source.")
	cmd.Flags().BoolVarP(&startFlags.IsTheme, "theme", "t", false, "Run the site as a theme using the current folder as the theme source.")
	cmd.Flags().BoolVarP(&startFlags.SSL, "ssl", "s", false, "Whether the site should default to SSL (https) or not.")
	cmd.Flags().BoolVarP(
		&startFlags.Activate,
		"activate",
		"a",
		false,
		"Activate the current plugin or theme (only works when used with the 'plugin' or 'theme' flags).")
	cmd.Flags().BoolVarP(
		&startFlags.Local,
		"local",
		"l",
		false,
		"Installs the WordPress files in your current path at ./wordpress instead of the global app path.")
	cmd.Flags().StringVarP(&startFlags.Multisite, "multisite", "u", "none", "Creates your new site as a multisite installation.")
	cmd.Flags().Lookup("network").NoOptDefVal = "subdomain"

	return cmd
}
