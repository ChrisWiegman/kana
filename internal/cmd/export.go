package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
)

func export(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the current config to a .kana.json file to save with your repo.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			if !kanaSite.IsSiteRunning() {
				consoleOutput.Error(fmt.Errorf("the export command only works on a running site.  Please run 'kana start' to start the site"))
			}

			err = kanaSite.ExportSiteConfig(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			consoleOutput.Success(fmt.Sprintf("Your config has been exported to %s", filepath.Join(kanaSite.Settings.WorkingDirectory, ".kana.json")))
		},
		Args: cobra.ArbitraryArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.DisableFlagParsing = true

	return cmd
}
