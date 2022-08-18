package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
)

func newDestroyCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroys the current WordPress site. This is a permanent change.",
		Run: func(cmd *cobra.Command, args []string) {
			runDestroy(cmd, args, site)
		},
		Args: cobra.NoArgs,
	}

	return cmd

}

func runDestroy(cmd *cobra.Command, args []string, site *site.Site) {

	err := site.StopWordPress()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.RemoveAll(site.StaticConfig.SiteDirectory)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
