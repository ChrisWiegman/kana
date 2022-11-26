package cmd

import (
	"github.com/ChrisWiegman/kana-cli/internal/config"
	"github.com/ChrisWiegman/kana-cli/internal/console"

	"github.com/spf13/cobra"
)

func newConfigCommand(kanaConfig *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and edit the saved configuration for the app or the local site.",
		Run: func(cmd *cobra.Command, args []string) {
			runConfigCommand(cmd, args, kanaConfig)
		},
		Args: cobra.RangeArgs(0, 2),
	}

	return cmd
}

func runConfigCommand(cmd *cobra.Command, args []string, kanaConfig *config.Config) {

	// List all content if we don't have args, list the value with 1 arg or set a fresh value with 2 args.
	// This is similar to how setting git options works
	switch len(args) {
	case 0:
		kanaConfig.ListDynamicContent()
	case 1:
		value, err := kanaConfig.GetDynamicContentItem(cmd, args)
		if err != nil {
			console.Error(err, flagVerbose)
		}

		console.Println(value)
	case 2:
		err := kanaConfig.SetDynamicContent(cmd, args)
		if err != nil {
			console.Error(err, flagVerbose)
		}
	}
}
