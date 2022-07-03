package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/config"

	"github.com/spf13/cobra"
)

func newRootCommand() (*cobra.Command, config.AppConfig) {

	appConfig, err := config.GetAppConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   "kana",
		Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
	}

	return cmd, appConfig

}

func Execute() {

	cmd, appConfig := newRootCommand()

	cmd.AddCommand(
		newStartCommand(appConfig),
		newStopCommand(appConfig),
		newOpenCommand(appConfig),
		newWPCommand(appConfig),
	)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
