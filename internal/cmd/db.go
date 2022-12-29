package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/spf13/cobra"
)

var flagPreserve bool
var flagReplaceDomain string

func newDBCommand(kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Commands to easily import and export a WordPress database from an existing site",
		Args:  cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	importCmd := &cobra.Command{
		Use:   "import <sql file>",
		Short: "Import a database from an existing WordPress site",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			err = kanaSite.ImportDatabase(args[0], flagPreserve, flagReplaceDomain)
			if err != nil {
				console.Error(err, flagVerbose)
			}

			console.Success("Your database file has been successfully imported and processed. Reload your site to see the changes.")
		},
		Args: cobra.ExactArgs(1),
	}

	commandsRequiringSite = append(commandsRequiringSite, importCmd.Use)

	exportCmd := &cobra.Command{
		Use:   "export [sql file]",
		Short: "Export the site's WordPress database",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := kanaSite.ExportDatabase(args)
			if err != nil {
				console.Error(err, flagVerbose)
			}

			console.Success(fmt.Sprintf("Export complete. Your database has been exported to %s.", file))
		},
		Args: cobra.MaximumNArgs(1),
	}

	commandsRequiringSite = append(commandsRequiringSite, exportCmd.Use)

	importCmd.Flags().BoolVarP(&flagPreserve, "preserve", "p", false, "Preserve the existing database (don't drop it before import)")
	importCmd.Flags().StringVarP(
		&flagReplaceDomain,
		"replace-domain",
		"d",
		"",
		"The old site domain to replace automatically with the development site domain")

	cmd.AddCommand(
		importCmd,
		exportCmd,
	)

	return cmd
}
