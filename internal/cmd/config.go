package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana-cli/internal/appConfig"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

func newConfigCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and edit the saved configuration for the app or the local site.",
		Run: func(cmd *cobra.Command, args []string) {
			runConfigCommand(cmd, args, site)
		},
		Args: cobra.RangeArgs(0, 2),
	}

	return cmd
}

func runConfigCommand(cmd *cobra.Command, args []string, site *site.Site) {

	// List all content if we don't have args, list the value with 1 arg or set a fresh value with 2 args.
	// This is similar to how setting git options works
	switch len(args) {
	case 0:
		appConfig.ListDynamicContent(site.DynamicConfig)
	case 1:
		value, err := appConfig.GetDynamicContentItem(cmd, args, site.DynamicConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(value)
	case 2:
		err := appConfig.SetDynamicContent(cmd, args, site.DynamicConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
