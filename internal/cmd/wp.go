package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

func newWPCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "wp",
		Short: "Run a wp-cli command against the current site.",
		Run: func(cmd *cobra.Command, args []string) {
			runWP(cmd, args, site)
		},
		Args: cobra.ArbitraryArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.DisableFlagParsing = true

	return cmd
}

func runWP(cmd *cobra.Command, args []string, site *site.Site) {

	if !site.IsSiteRunning() {
		console.Error(fmt.Errorf("the `wp` command only works on a running site. Please run 'kana start' to start the site"), flagDebugMode)
	}

	// Run the output from wp-cli
	output, err := site.RunWPCli(args)
	if err != nil {
		console.Error(err, flagDebugMode)
	}

	console.Println(output)
}
