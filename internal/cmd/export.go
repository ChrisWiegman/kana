package cmd

import (
	"fmt"
	"path"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/spf13/cobra"
)

func newExportCommand(kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the current config to a .kana.json file to save with your repo.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			if !kanaSite.IsSiteRunning() {
				console.Error(fmt.Errorf("the export command only works on a running site.  Please run 'kana start' to start the site"), flagVerbose)
			}

			err = kanaSite.ExportSiteConfig()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			console.Success(fmt.Sprintf("Your config has been exported to %s", path.Join(kanaSite.Settings.WorkingDirectory, ".kana.json")))
		},
		Args: cobra.ArbitraryArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.DisableFlagParsing = true

	return cmd
}
