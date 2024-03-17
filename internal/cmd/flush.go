package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

func flush(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flush",
		Short: "Flushes the cache and deletes all transients.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			commands := [][]string{
				{"cache", "flush"},
				{"transient", "delete", "--all"},
			}

			for _, command := range commands {
				var code int64

				code, _, err = kanaSite.RunWPCli(command, false, consoleOutput)
				if code != 0 || err != nil {
					consoleOutput.Error(fmt.Errorf("unable to complete flush"))
				}
			}

			consoleOutput.Success("Cache and transients have been successfully flushed")
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}
