package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/aquasecurity/table"
	"github.com/spf13/cobra"
)

func newListCommand(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Kana sites and their associated status.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			sites, err := kanaSite.GetSiteList(kanaSite.Settings.AppDirectory, consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			if consoleOutput.JSON {
				str, _ := json.Marshal(sites)

				fmt.Println(string(str))

				return
			}

			t := table.New(os.Stdout)

			t.SetHeaders("Name", "Path", "Running")

			for _, site := range sites {
				t.AddRow(site.Name, site.Path, strconv.FormatBool(site.Running))
			}

			t.Render()

			if len(sites) == 0 {
				consoleOutput.Println("You have not created any sites with Kana yet. Use `kana start` to create a site.")
			}
		},
		Args: cobra.NoArgs,
	}

	return cmd
}
