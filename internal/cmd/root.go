package cmd

import (
	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

var (
	flagName                    string
	flagVerbose, flagJSONOutput bool
	commandsRequiringSite       []string
)

func Execute() {
	kanaSite := new(site.Site)
	consoleOutput := new(console.Console)

	// Setup the cobra command
	cmd := &cobra.Command{
		Use:   "kana",
		Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
		Args:  cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			consoleOutput.Debug = flagVerbose
			consoleOutput.JSON = flagJSONOutput

			err := kanaSite.LoadSite(cmd, commandsRequiringSite, startFlags, flagVerbose)
			if err != nil {
				consoleOutput.Error(err)
			}
		},
	}

	// Hide the default completion command
	cmd.CompletionOptions.HiddenDefaultCmd = true

	// Add the "name" flag to allow for sites not connected to the local directory
	cmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "Specify a name for the site, used to override using the current folder.")
	cmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Display debugging information along with detailed command output")
	cmd.PersistentFlags().BoolVar(&flagJSONOutput, "output-json", false, "Display all output in JSON format for further processing")

	err := cmd.PersistentFlags().MarkHidden("output-json")
	if err != nil {
		consoleOutput.Error(err)
	}

	// Register the subcommands
	cmd.AddCommand(
		newStartCommand(consoleOutput, kanaSite),
		newStopCommand(consoleOutput, kanaSite),
		newOpenCommand(consoleOutput, kanaSite),
		newWPCommand(consoleOutput, kanaSite),
		newDestroyCommand(consoleOutput, kanaSite),
		newConfigCommand(consoleOutput, kanaSite),
		newExportCommand(consoleOutput, kanaSite),
		newVersionCommand(consoleOutput),
		newDBCommand(consoleOutput, kanaSite),
		newListCommand(consoleOutput, kanaSite),
		newXdebugCommand(consoleOutput, kanaSite),
		newFlushCommand(consoleOutput, kanaSite),
	)

	// Execute anything we need to
	if err := cmd.Execute(); err != nil {
		consoleOutput.Error(err)
	}
}
