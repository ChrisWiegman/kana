package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/options"
)

func test(consoleOutput *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "A test command for developing a simpler settings package.",
		Run: func(cmd *cobra.Command, args []string) {
			kanaSettings, err := options.New(Version, cmd)
			if err != nil {
				panic(err)
			}

			options.ListSettings(kanaSettings, consoleOutput)
		},
		Args: cobra.NoArgs,
	}

	return cmd
}
