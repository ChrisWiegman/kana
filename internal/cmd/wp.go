package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/app"
	"github.com/ChrisWiegman/kana-cli/internal/config"
	"github.com/ChrisWiegman/kana-cli/internal/console"

	"github.com/spf13/cobra"
)

func newWPCommand(kanaConfig *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "wp",
		Short: "Run a wp-cli command against the current site.",
		Run: func(cmd *cobra.Command, args []string) {
			runWP(cmd, args, kanaConfig)
		},
		Args: cobra.ArbitraryArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.DisableFlagParsing = true

	return cmd
}

func runWP(cmd *cobra.Command, args []string, kanaConfig *config.Config) {

	site, err := app.NewSite(kanaConfig)
	if err != nil {
		console.Error(err, flagVerbose)
	}

	if !site.IsSiteRunning() {
		console.Error(fmt.Errorf("the `wp` command only works on a running site. Please run 'kana start' to start the site"), flagVerbose)
	}

	// Run the output from wp-cli
	output, err := site.RunWPCli(args)
	if err != nil {
		console.Error(err, flagVerbose)
	}

	console.Println(output)
}
