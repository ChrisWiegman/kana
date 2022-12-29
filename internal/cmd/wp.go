package cmd

import (
	"errors"
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/site"
	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/spf13/cobra"
)

func newWPCommand(kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wp",
		Short: "Run a wp-cli command against the current site.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker()
			if err != nil {
				console.Error(err, flagVerbose)
			}

			if !kanaSite.IsSiteRunning() {
				console.Error(fmt.Errorf("the `wp` command only works on a running site. Please run 'kana start' to start the site"), flagVerbose)
			}

			// Run the output from wp-cli
			code, output, err := kanaSite.RunWPCli(args)
			if err != nil {
				console.Error(err, flagVerbose)
			}

			if code != 0 {
				console.Error(errors.New(output), flagVerbose)
			}

			console.Println(output)
		},
		Args: cobra.ArbitraryArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.DisableFlagParsing = true

	return cmd
}
