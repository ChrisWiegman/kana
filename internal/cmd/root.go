package cmd

import (
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/spf13/cobra"
)

var flagName string
var flagVerbose, flagJsonOutput bool
var commandsRequiringSite []string

func Execute() {
	kanaSite := new(site.Site)

	// Setup the cobra command
	cmd := &cobra.Command{
		Use:   "kana",
		Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
		Args:  cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := kanaSite.LoadSite(cmd, commandsRequiringSite, startFlags, flagVerbose)
			if err != nil {
				console.Error(err, flagVerbose)
			}
		},
	}

	// Hide the default completion command
	cmd.CompletionOptions.HiddenDefaultCmd = true

	// Add the "name" flag to allow for sites not connected to the local directory
	cmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "Specify a name for the site, used to override using the current folder.")
	cmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Display debugging information along with detailed command output")
	cmd.PersistentFlags().BoolVar(&flagJsonOutput, "output-json", false, "Display all output in JSON format for further processing")

	err := cmd.PersistentFlags().MarkHidden("output-json")
	if err != nil {
		console.Error(err, flagVerbose)
	}

	// Register the subcommands
	cmd.AddCommand(
		newStartCommand(kanaSite),
		newStopCommand(kanaSite),
		newOpenCommand(kanaSite),
		newWPCommand(kanaSite),
		newDestroyCommand(kanaSite),
		newConfigCommand(kanaSite),
		newExportCommand(kanaSite),
		newVersionCommand(),
		newDBCommand(kanaSite),
		newListCommand(kanaSite),
	)

	// Execute anything we need to
	if err := cmd.Execute(); err != nil {
		console.Error(err, flagVerbose)
	}
}
