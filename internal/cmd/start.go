package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/config"
	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/internal/traefik"

	"github.com/spf13/cobra"
)

var flagXdebug bool
var flagLocal bool
var flagIsTheme bool
var flagIsPlugin bool

func newStartCommand(kanaConfig *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new environment in the local folder.",
		Run: func(cmd *cobra.Command, args []string) {

			site, err := site.NewSite(kanaConfig)
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// A site shouldn't be both a plugin and a theme so this reports an error if that is the case.
			if flagIsPlugin && flagIsTheme {
				console.Error(fmt.Errorf("you have set both the plugin and theme flags. Please choose only one option"), flagVerbose)
			}

			// Check that the site is already running and show an error if it is.
			if site.IsSiteRunning() {
				console.Error(fmt.Errorf("the site is already running. Please stop your site before running the start command"), flagVerbose)
			}

			// Process any overrides set with flags on the start command
			startFlags := config.StartFlags{
				Xdebug:   flagXdebug,
				IsTheme:  flagIsTheme,
				IsPlugin: flagIsPlugin,
				Local:    flagLocal,
			}

			kanaConfig.ProcessStartFlags(cmd, startFlags)

			// Let's start everything up
			fmt.Printf("Starting development site: %s\n", site.GetURL(false))

			// Start Traefik if we need it
			traefikClient, err := traefik.NewTraefik(kanaConfig)
			if err != nil {
				console.Error(err, flagVerbose)
			}

			err = traefikClient.StartTraefik()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// Start WordPress
			err = site.StartWordPress()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// Make sure the WordPress site is running
			_, err = site.VerifySite()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// Setup WordPress
			err = site.InstallWordPress()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// Install Xdebug if we need to
			_, err = site.InstallXdebug()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// Install any configuration plugins if needed
			err = site.InstallDefaultPlugins()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			// Open the site in the user's browser
			err = site.OpenSite()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			console.Success("Your site has started and should be open in your default browser.")
		},
		Args: cobra.NoArgs,
	}

	// Add associated flags to customize the site at runtime.
	cmd.Flags().BoolVarP(&flagXdebug, "xdebug", "x", false, "Enable Xdebug when starting the container.")
	cmd.Flags().BoolVarP(&flagIsPlugin, "plugin", "p", false, "Run the site as a plugin using the current folder as the plugin source.")
	cmd.Flags().BoolVarP(&flagIsTheme, "theme", "t", false, "Run the site as a theme using the current folder as the theme source.")
	cmd.Flags().BoolVarP(&flagLocal, "local", "l", false, "Installs the WordPress files in your current path at ./wordpress instead of the global app path.")

	return cmd
}
