package cmd

import (
	"os"

	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/aquasecurity/table"
	"github.com/spf13/cobra"
)

func newConfigCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Edit the saved configuration for the app or the local site.",
	}

	cmd.AddCommand(
		newConfigListCommand(site),
	)

	return cmd

}

func newConfigListCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all config items and their values.",
		Run: func(cmd *cobra.Command, args []string) {
			runConfig(cmd, args, site)
		},
		Args: cobra.ExactArgs(0),
	}

	return cmd
}

func runConfig(cmd *cobra.Command, args []string, site *site.Site) {

	t := table.New(os.Stdout)

	t.SetHeaders("Key", "Value")

	t.AddRow("adminEmail", site.DynamicConfig.GetString("adminEmail"))
	t.AddRow("adminPassword", site.DynamicConfig.GetString("adminPassword"))
	t.AddRow("adminUser", site.DynamicConfig.GetString("adminUser"))
	t.AddRow("local", site.DynamicConfig.GetString("local"))
	t.AddRow("php", site.DynamicConfig.GetString("php"))
	t.AddRow("type", site.DynamicConfig.GetString("type"))
	t.AddRow("xdebug", site.DynamicConfig.GetString("xdebug"))

	t.Render()

}
