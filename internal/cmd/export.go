package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

func newExportCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the current config to a .kana.json file to save with your repo.",
		Run: func(cmd *cobra.Command, args []string) {
			runExport(cmd, args, site)
		},
		Args: cobra.ArbitraryArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.DisableFlagParsing = true

	return cmd
}

func runExport(cmd *cobra.Command, args []string, site *site.Site) {

	if !site.IsSiteRunning() {
		console.Error(fmt.Errorf("the export command only works on a running site.  Please run 'kana start' to start the site"), flagDebugMode)
	}

	err := site.ExportSiteConfig()
	if err != nil {
		console.Error(err, flagDebugMode)
	}
}
