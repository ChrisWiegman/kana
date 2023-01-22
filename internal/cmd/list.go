package cmd

import (
	"os"
	"strconv"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/aquasecurity/table"
	"github.com/spf13/cobra"
)

func newListCommand(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Kana sites and their associated status.",
		Run: func(cmd *cobra.Command, args []string) {
			sites, err := site.GetSiteList(kanaSite.Settings.AppDirectory, consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			if len(sites) > 0 {
				t := table.New(os.Stdout)

				t.SetHeaders("Name", "Path", "Running")

				for _, site := range sites {
					t.AddRow(site.Name, site.Path, strconv.FormatBool(site.Running))
				}

				t.Render()
			} else {
				consoleOutput.Println("It doesn't look like you have created any sites with Kana yet. Use `kana start` to create a site.")
			}
		},
		Args: cobra.NoArgs,
	}

	return cmd
}
