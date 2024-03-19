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

func list(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Kana sites and their associated status.",
		Run: func(cmd *cobra.Command, args []string) {
			dockerIsRunning := true

			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				dockerIsRunning = false
			}

			sites, err := kanaSite.GetSiteList(kanaSite.Settings.AppDirectory, dockerIsRunning)
			if err != nil {
				consoleOutput.Error(err)
			}

			if consoleOutput.JSON {
				if len(sites) > 0 {
					str, _ := json.Marshal(sites)

					fmt.Println(string(str))
				} else {
					fmt.Println("[]")
				}

				return
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
