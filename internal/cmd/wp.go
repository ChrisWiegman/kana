package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/setup"
	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
)

func newWPCommand(appConfig config.AppConfig) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "wp",
		Short: "Run a wp-cli command against the current site.",
		Run: func(cmd *cobra.Command, args []string) {
			runWP(cmd, args, appConfig)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setup.SetupApp(appConfig)
		},
	}

	cmd.DisableFlagParsing = true

	return cmd

}

func runWP(cmd *cobra.Command, args []string, appConfig config.AppConfig) {

	site, err := site.NewSite(appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output, err := site.RunWPCli(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(output)
}
