package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/setup"

	"github.com/spf13/cobra"
)

func newConfigCommand(appConfig config.AppConfig) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Edit the saved configuration for the app or the local site.",
		Run: func(cmd *cobra.Command, args []string) {
			runConfig(cmd, args, appConfig)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setup.SetupApp(appConfig)
		},
	}

	return cmd

}

func runConfig(cmd *cobra.Command, args []string, appConfig config.AppConfig) {

	fmt.Println(appConfig)
}
