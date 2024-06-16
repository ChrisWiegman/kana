package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/site"

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

			if cmd.Flags().Lookup("name").Changed {
				for i := range args {
					if !strings.Contains(args[i], "--name=") && !strings.Contains(args[i], "-n=") {
						continue
					}

					nameFlag := strings.Split(args[i], "=")

					if nameFlag[1] == cmd.Flags().Lookup("name").Value.String() {
						args = append(args[:i], args[i+1:]...)
						break
					}
				}
			}

			// Run the output from wp-cli
			code, output, err := kanaSite.Cli.WPCli(args, true, consoleOutput)
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

func parseWPNameFlag(args []string, cmd *cobra.Command) error {
	for i := range args {
		if !strings.Contains(args[i], "--name=") && !strings.Contains(args[i], "-n=") {
			continue
		}

		nameFlag := strings.Split(args[i], "=")

		err := cmd.Flag("name").Value.Set(nameFlag[1])
		if err != nil {
			return err
		}

		cmd.Flag("name").Changed = true
		break
	}

	return nil
}
