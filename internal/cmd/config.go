package cmd

import (
	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/settings"

	"github.com/spf13/cobra"
)

func config(consoleOutput *console.Console, kanaSettings *settings.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and edit the saved configuration for the app or the local site.",
		Run: func(cmd *cobra.Command, args []string) {
			// List all content if we don't have args, list the value with 1 arg or set a fresh value with 2 args.
			// This is similar to how setting git options works
			switch len(args) {
			case 0:
				settings.ListSettings(kanaSettings, consoleOutput)
			case 1:
				kanaSettings.PrintSingleSetting(args[0], consoleOutput)
			case 2:
				err := kanaSettings.Set(args[0], args[1])
				if err != nil {
					consoleOutput.Error(err)
				}

				kanaSettings.PrintSingleSetting(args[0], consoleOutput)
			}
		},
		Args: cobra.RangeArgs(0, 2),
	}

	return cmd
}
