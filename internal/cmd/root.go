package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/settings"
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/spf13/cobra"
)

var flagName string
var flagVerbose bool
var commandsRequiringSite []string

func Execute() {

	site, err := site.NewSite()
	if err != nil {
		console.Error(err, flagVerbose)
	}

	// Setup the cobra command
	cmd := &cobra.Command{
		Use:   "kana",
		Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
		Args:  cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			// Process the "name" flag for every command
			isSite, err := site.Settings.ProcessNameFlag(cmd)
			if err != nil {
				console.Error(err, flagVerbose)
			}

			if !isSite && arrayContains(commandsRequiringSite, cmd.Use) {
				console.Error(fmt.Errorf("the current site you are trying to work with does not exist. Use `kana start` to initialize"), flagVerbose)
			}

			// Process the "start" command flags
			if cmd.Use == "start" {

				// A site shouldn't be both a plugin and a theme so this reports an error if that is the case.
				if flagIsPlugin && flagIsTheme {
					console.Error(fmt.Errorf("you have set both the plugin and theme flags. Please choose only one option"), flagVerbose)
				}

				// Process any overrides set with flags on the start command
				startFlags := settings.StartFlags{
					Xdebug:   flagXdebug,
					IsTheme:  flagIsTheme,
					IsPlugin: flagIsPlugin,
					Local:    flagLocal,
				}

				site.Settings.ProcessStartFlags(cmd, startFlags)
			}
		},
	}

	// Add the "name" flag to allow for sites not connected to the local directory
	cmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "Specify a name for the site, used to override using the current folder.")
	cmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Display debugging information along with detailed command output")

	// Register the subcommands
	cmd.AddCommand(
		newStartCommand(site),
		newStopCommand(site),
		newOpenCommand(site),
		newWPCommand(site),
		newDestroyCommand(site),
		newConfigCommand(site),
		newExportCommand(site),
		newVersionCommand(),
		newDbCommand(site),
	)

	// Execute anything we need to
	if err := cmd.Execute(); err != nil {
		console.Error(err, flagVerbose)
	}
}

func arrayContains(array []string, name string) bool {
	for _, value := range array {
		if value == name {
			return true
		}
	}

	return false
}
