package cmd

import (
	"os"
	"strconv"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/aquasecurity/table"

	"github.com/spf13/cobra"
)

func newListCommand(kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Kana sites and their associated status.",
		Run: func(cmd *cobra.Command, args []string) {
			sites, err := site.GetSiteList(kanaSite.Settings.AppDirectory)
			if err != nil {
				console.Error(err, flagVerbose)
			}

			t := table.New(os.Stdout)

			t.SetHeaders("Name", "Path", "Running")

			for _, site := range sites {
				t.AddRow(site.Name, site.Path, strconv.FormatBool(site.Running))
			}

			t.Render()
		},
		Args: cobra.NoArgs,
	}

	return cmd
}
