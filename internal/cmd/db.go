package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana-cli/internal/database"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

var flagPreserve bool
var flagReplaceDomain string

func newDbCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "db",
		Short: "Commands to easily import and export a WordPress database from an existing site",
		Args:  cobra.NoArgs,
	}

	importCmd := &cobra.Command{
		Use:   "import <sql file>",
		Short: "Import a database from an existing WordPress site",
		Run: func(cmd *cobra.Command, args []string) {
			runDbImport(cmd, args, site)
		},
		Args: cobra.ExactArgs(1),
	}

	exportCmd := &cobra.Command{
		Use:   "export [sql file]",
		Short: "Export the site's WordPress database",
		Run: func(cmd *cobra.Command, args []string) {
			runDbExport(cmd, args, site)
		},
		Args: cobra.MaximumNArgs(1),
	}

	importCmd.Flags().BoolVarP(&flagPreserve, "preserve", "p", false, "Preserve the existing database (don't drop it before import)")
	importCmd.Flags().StringVarP(&flagReplaceDomain, "replace-domain", "d", "", "The old site domain to replace automatically with the development site domain")

	cmd.AddCommand(
		importCmd,
		exportCmd,
	)

	return cmd
}

func runDbImport(cmd *cobra.Command, args []string, kanaSite *site.Site) {
	err := database.Import(kanaSite, args[0], flagPreserve, flagReplaceDomain)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runDbExport(cmd *cobra.Command, args []string, kanaSite *site.Site) {
	err := database.Export(kanaSite, args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
