package cmd

import (
	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/settings"
	"github.com/ChrisWiegman/kana-cli/internal/site"

	"github.com/spf13/cobra"
)

func trust(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trust-ssl",
		Short: "Add the Kana SSL certificate to the MacOS Keychain (if needed).",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.Settings.EnsureSSLCerts(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			err = settings.TrustSSL(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			consoleOutput.Success("The Kana SSL certificate has been added to your Mac's system Keychain.")
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)

	return cmd
}
