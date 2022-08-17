package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
)

func newConfigCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Edit the saved configuration for the app or the local site.",
		Run: func(cmd *cobra.Command, args []string) {
			runConfig(cmd, args, site)
		},
	}

	return cmd

}

func runConfig(cmd *cobra.Command, args []string, site *site.Site) {
	fmt.Println(site.StaticConfig)
}
