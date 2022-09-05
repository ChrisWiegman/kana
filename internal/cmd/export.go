package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/site"

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

	cmd.DisableFlagParsing = true

	return cmd
}

func runExport(cmd *cobra.Command, args []string, site *site.Site) {

	if !site.IsSiteRunning() {
		fmt.Println("The export command only works on a running site.  Please run 'kana start' to start the site.")
		os.Exit(1)
	}

	err := site.ExportSiteSettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
