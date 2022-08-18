package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/appConfig"
	"github.com/ChrisWiegman/kana/internal/site"

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

	switch len(args) {
	case 0:
		appConfig.ListDynamicContent(site.DynamicConfig)
	case 1:
		fmt.Println("get the config item")
	case 2:
		err := appConfig.SetDynamicContent(cmd, args, site.DynamicConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}
