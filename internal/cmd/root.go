package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana-cli/internal/appConfig"
	"github.com/ChrisWiegman/kana-cli/internal/appSetup"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

var flagName string

func Execute() {

	// Setup the static config items that cannot be overripen
	staticConfig, err := appConfig.GetStaticConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Ensure the static content files are in place and up to date
	err = appSetup.EnsureStaticConfigFiles(staticConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Get the dynamic config that the user might have set themselves
	dynamicConfig, err := appConfig.GetDynamicContent(staticConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Create a site object
	site, err := site.NewSite(staticConfig, dynamicConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Setup the cobra command
	cmd := &cobra.Command{
		Use:   "kana",
		Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
		Args:  cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := site.ProcessNameFlag(cmd)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	// Add the "name" flag to allow for sites not connected to the local directory
	cmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "Specify a name for the site, used to override using the current folder.")

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
	)

	// Execute anything we need to
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
