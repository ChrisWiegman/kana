package cmd

import (
	"github.com/ChrisWiegman/kana-wordpress/internal/console"
	"github.com/ChrisWiegman/kana-wordpress/internal/settings"

	"github.com/spf13/cobra"
)

func trust(consoleOutput *console.Console, kanaSettings *settings.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trust-ssl",
		Short: "Add the Kana SSL certificate to the MacOS Keychain (if needed).",
		Run: func(cmd *cobra.Command, args []string) {
			err := settings.EnsureSSLCerts(
				kanaSettings.Get("appDirectory"),
				kanaSettings.GetBool("ssl"),
				consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			err = settings.TrustSSL(kanaSettings.Get("rootCert"), consoleOutput)
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
