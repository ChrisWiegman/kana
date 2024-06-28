package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/options"
)

func test(consoleOutput *console.Console, kanaSettings *options.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "A test command for developing a simpler settings package.",
		Run: func(cmd *cobra.Command, args []string) {
			// List all content if we don't have args, list the value with 1 arg or set a fresh value with 2 args.
			// This is similar to how setting git options works
			switch len(args) {
			case 0:
				options.ListSettings(kanaSettings, consoleOutput)
			case 1:
				options.PrintSingleSetting(args[0], kanaSettings, consoleOutput)
			case 2:
				err := kanaSettings.Set(args[0], args[1], true)
				if err != nil {
					consoleOutput.Error(err)
				}

				options.PrintSingleSetting(args[0], kanaSettings, consoleOutput)
			}
		},
		Args: cobra.RangeArgs(0, 2),
	}

	options.AddStartFlags(cmd, kanaSettings)

	return cmd
}
