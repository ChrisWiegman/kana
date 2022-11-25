package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/appConfig"
	"github.com/ChrisWiegman/kana-cli/internal/appSetup"
	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

var flagName string
var flagDebugMode bool
var commandsRequiringSite []string

func Execute() {

	// Setup the static config items that cannot be overripen
	staticConfig, err := appConfig.GetStaticConfig()
	if err != nil {
		console.Error(err, flagDebugMode)
	}

	// Ensure the static content files are in place and up to date
	err = appSetup.EnsureStaticConfigFiles(staticConfig)
	if err != nil {
		console.Error(err, flagDebugMode)
	}

	// Get the dynamic config that the user might have set themselves
	dynamicConfig, err := appConfig.GetDynamicContent(staticConfig)
	if err != nil {
		console.Error(err, flagDebugMode)
	}

	// Create a site object
	site, err := site.NewSite(staticConfig, dynamicConfig)
	if err != nil {
		console.Error(err, flagDebugMode)
	}

	// Setup the cobra command
	cmd := &cobra.Command{
		Use:   "kana",
		Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
		Args:  cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			isSite, err := site.ProcessNameFlag(cmd)
			if err != nil {
				console.Error(err, flagDebugMode)
			}

			if !isSite && arrayContains(commandsRequiringSite, cmd.Use) {
				console.Error(fmt.Errorf("the current site you are trying to work with does not exist. Use `kana start` to initialize"), flagDebugMode)
			}
		},
	}

	// Add the "name" flag to allow for sites not connected to the local directory
	cmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "Specify a name for the site, used to override using the current folder.")
	cmd.PersistentFlags().BoolVarP(&flagDebugMode, "debug", "d", false, "Display debugging information along with the command output")

	// Register the subcommands
	cmd.AddCommand(
		newStartCommand(site),
		newStopCommand(site),
		newOpenCommand(site),
		newWPCommand(site),
		newDestroyCommand(site),
		newConfigCommand(site),
		newExportCommand(site),
		newVersionCommand(site),
		newDbCommand(site),
	)

	// Execute anything we need to
	if err := cmd.Execute(); err != nil {
		console.Error(err, flagDebugMode)
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
