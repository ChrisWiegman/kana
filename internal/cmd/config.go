package cmd

import (
	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/spf13/cobra"
)

func newConfigCommand(kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and edit the saved configuration for the app or the local site.",
		Run: func(cmd *cobra.Command, args []string) {
			// List all content if we don't have args, list the value with 1 arg or set a fresh value with 2 args.
			// This is similar to how setting git options works
			switch len(args) {
			case 0:
				kanaSite.Settings.ListSettings()
			case 1:
				value, err := kanaSite.Settings.GetGlobalSetting(cmd, args)
				if err != nil {
					console.Error(err, flagVerbose)
				}

				console.Println(value)
			case 2:
				err := kanaSite.Settings.SetGlobalSetting(cmd, args)
				if err != nil {
					console.Error(err, flagVerbose)
				}
			}
		},
		Args: cobra.RangeArgs(0, 2),
	}

	return cmd
}
