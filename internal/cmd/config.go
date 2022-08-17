package cmd

import (
	"github.com/ChrisWiegman/kana/internal/appConfig"
	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
)

func newConfigCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Edit the saved configuration for the app or the local site.",
	}

	cmd.AddCommand(
		newConfigListCommand(site),
		newConfigSetCommand(site),
	)

	return cmd

}

func newConfigListCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all config items and their values.",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig.ListDynamicContent(site.DynamicConfig)
		},
		Args: cobra.ExactArgs(0),
	}

	return cmd
}

func newConfigSetCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set a new config value.",
		Run: func(cmd *cobra.Command, args []string) {
			appConfig.SetDynamicContent(cmd, args, site.DynamicConfig)
		},
		Args: cobra.ExactArgs(2),
	}

	return cmd
}
