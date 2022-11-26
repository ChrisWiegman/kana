package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/config"
	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

func newWPCommand(kanaConfig *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "wp",
		Short: "Run a wp-cli command against the current site.",
		Run: func(cmd *cobra.Command, args []string) {

			site, err := site.NewSite(kanaConfig)
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
		},
		Args: cobra.ArbitraryArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.DisableFlagParsing = true

	return cmd
}
