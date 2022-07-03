package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/setup"
	"github.com/ChrisWiegman/kana/internal/wordpress"

	"github.com/spf13/cobra"
)

func newWPCommand(controller *docker.Controller) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "wp",
		Short: "Run a wp-cli command against the current site.",
		Run: func(cmd *cobra.Command, args []string) {
			runWP(cmd, args, controller)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setup.SetupApp(controller)
		},
	}

	cmd.DisableFlagParsing = true

	return cmd

}

func runWP(cmd *cobra.Command, args []string, controller *docker.Controller) {

	err := wordpress.RunCli(args, controller)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
