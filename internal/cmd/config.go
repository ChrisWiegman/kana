package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

func config(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and edit the saved configuration for the app or the local site.",
		Run: func(cmd *cobra.Command, args []string) {
			// List all content if we don't have args, list the value with 1 arg or set a fresh value with 2 args.
			// This is similar to how setting git options works
			switch len(args) {
			case 0:
				kanaSite.Settings.ListSettings(consoleOutput)
			case 1:
				value, err := kanaSite.Settings.GetGlobalSetting(args)
				if err != nil {
					consoleOutput.Error(err)
				}
				if consoleOutput.JSON {
					type JSONSetting struct {
						Setting, Value string
					}

					setting := JSONSetting{
						Setting: args[0],
						Value:   value,
					}

					str, _ := json.Marshal(setting)

					fmt.Println(string(str))
				} else {
					consoleOutput.Println(value)
				}
			case 2:
				err := kanaSite.Settings.SetGlobalSetting(args)
				if err != nil {
					consoleOutput.Error(err)
				}
			}
		},
		Args: cobra.RangeArgs(0, 2),
	}

	return cmd
}
