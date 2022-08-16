package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
)

func newDestroyCommand(appConfig config.AppConfig) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroys the current WordPress site. This is a permanent change.",
		Run: func(cmd *cobra.Command, args []string) {
			runDestroy(cmd, args, appConfig)
		},
	}

	return cmd

}

func runDestroy(cmd *cobra.Command, args []string, appConfig config.AppConfig) {

	site, err := site.NewSite(appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = site.StopWordPress()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.RemoveAll(appConfig.SiteDirectory)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
