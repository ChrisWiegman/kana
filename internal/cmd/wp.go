package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/site"

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

	cmd.DisableFlagParsing = true

	return cmd
}

func runWP(cmd *cobra.Command, args []string, site *site.Site) {

	if !site.IsSiteRunning() {
		fmt.Println("The wp command only works on a running site. Please run 'kana start' to start the site.")
		os.Exit(1)
	}

	// Run the output from wp-cli
	output, err := site.RunWPCli(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(output)
}
