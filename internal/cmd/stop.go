package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
)

func newStopCommand(appConfig config.AppConfig) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops the WordPress development environment.",
		Run: func(cmd *cobra.Command, args []string) {
			runStop(cmd, args, appConfig)
		},
	}

	return cmd

}

func runStop(cmd *cobra.Command, args []string, appConfig config.AppConfig) {

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
}
