package cmd

import (
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

func newListCommand(kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Kana sites and their associated status.",
		Run: func(cmd *cobra.Command, args []string) {
			site.GetSiteList(kanaSite.Settings.AppDirectory)
		},
		Args: cobra.NoArgs,
	}

	return cmd
}
