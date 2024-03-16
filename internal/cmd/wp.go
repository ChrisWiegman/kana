package cmd

import (
	"errors"
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

func wp(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wp",
		Short: "Run a wp-cli command against the current site.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			if !kanaSite.IsSiteRunning() {
				consoleOutput.Error(fmt.Errorf("the `wp` command only works on a running site. Please run 'kana start' to start the site"))
			}

			// Run the output from wp-cli
			code, output, err := kanaSite.RunWPCli(args, consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			if code != 0 {
				consoleOutput.Error(errors.New(output))
			}

			consoleOutput.Println(output)
		},
		Args: cobra.ArbitraryArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	cmd.DisableFlagParsing = true

	return cmd
}
