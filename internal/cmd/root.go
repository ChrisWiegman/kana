package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/setup"

	"github.com/spf13/cobra"
)

func Execute() {

	appConfig, err := config.GetAppConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = setup.SetupApp(appConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   "kana",
		Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
	}

	cmd.AddCommand(
		newStartCommand(appConfig),
		newStopCommand(appConfig),
		newOpenCommand(appConfig),
		newWPCommand(appConfig),
		newDestroyCommand(appConfig),
		newConfigCommand(appConfig),
	)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
