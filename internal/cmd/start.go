package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/helpers"
	"github.com/ChrisWiegman/kana/internal/settings"
	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func start(consoleOutput *console.Console, kanaSite *site.Site, kanaSettings *settings.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new environment in the local folder.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				if kanaSettings.GetBool("IsNew") {
					remError := os.RemoveAll(kanaSettings.Get("Site"))
					if remError != nil {
						consoleOutput.Error(remError)
					}
				}
				consoleOutput.Error(err)
			}

			err = handleTypeDetection(cmd, consoleOutput, kanaSettings)
			if err != nil {
				consoleOutput.Error(err)
			}

			if cmd.Flags().Lookup("theme").Changed && kanaSettings.Get("Type") == "theme" {
				consoleOutput.Error(fmt.Errorf("a default theme cannot be set on a site of type 'theme"))
			}

			// Check that the site is already running and show an error if it is.
			if kanaSite.IsSiteRunning() {
				consoleOutput.Error(fmt.Errorf("the site is already running. Please stop your site before running the start command"))
			}

			// Check that we're not using our home directory as the working directory as that could cause security or other issues.
			home, err := homedir.Dir()
			if err != nil {
				consoleOutput.Error(err)
			}

			if home == kanaSettings.Get("Working") {
				// Remove the site's folder in the config directory.
				err = os.RemoveAll(kanaSettings.Get("Site"))
				if err != nil {
					consoleOutput.Error(err)
				}

				consoleOutput.Error(fmt.Errorf("you are attempting to start a new site from your home directory. This could create security issues. Please create a folder and start a site from there")) //nolint:lll
			}

			err = kanaSite.StartSite(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			consoleOutput.Success(
				fmt.Sprintf(
					"Your site, %s, has has started and should be open in your default browser.",
					consoleOutput.Bold(consoleOutput.Blue(kanaSettings.Get("Name")))))
		},
		Args: cobra.NoArgs,
	}

	// Add associated flags to customize the site at runtime.
	cmd.Flags().BoolP("xdebug", "x", false, "Enable Xdebug when starting the container.")
	cmd.Flags().BoolP("scriptdebug", "c", false, "Enable SCRIPT_DEBUG when starting the container.")
	cmd.Flags().BoolP("wpdebug", "d", false, "Enable WP_Debug when starting the container.")
	cmd.Flags().BoolP("mailpit", "m", false, "Enable Mailpit when starting the container.")
	cmd.Flags().BoolP("ssl", "s", false, "Whether the site should default to SSL (https) or not.")
	cmd.Flags().BoolP("removedefaultplugins",
		"r",
		false,
		"If true will remove the default plugins installed with WordPress (Akismet and Hello Dolly) when starting a site.")
	cmd.Flags().BoolP("activate",
		"a",
		false,
		"Activate the current plugin or theme (only works when used with the 'plugin' or 'theme' flags).")

	cmd.Flags().String("database", "mariadb", "Select the database server you wish to use with your installation.")
	cmd.Flags().String("multisite", "none", "Creates your new site as a multisite installation.")
	cmd.Flags().String("environment", "local", "Sets the WP_ENVIRONMENT_TYPE for the site.")
	cmd.Flags().String("plugins",
		"",
		"Installs and activates the specified plugins. Multiple plugins should be separated by commas")
	cmd.Flags().String("theme", "", "Installs and activates a theme when starting a site.")
	cmd.Flags().String("type", "site", "Set the type of the installation, `site`, `plugin` or `theme`.")

	cmd.Flags().Lookup("multisite").NoOptDefVal = "subdomain"

	return cmd
}

func handleTypeDetection(cmd *cobra.Command, consoleOutput *console.Console, kanaSettings *settings.Settings) error {
	if !cmd.Flags().Lookup("type").Changed && !kanaSettings.GetBool("HasLocalSettings") {
		if !cmd.Flags().Lookup("name").Changed {
			err := verifyEmpty(kanaSettings, consoleOutput)
			if err != nil {
				return err
			}
		}

		if kanaSettings.GetBool("TypeIsDetected") {
			consoleOutput.Printf(
				"A %s was detected in the current site folder. Starting site as a %s\n",
				kanaSettings.Get("Type"),
				kanaSettings.Get("Type"))
		}
	}

	return nil
}

// verifyEmpty Verifies the folder is empty when starting a new site in it.
// This helps prevent conflicts with WordPress files and anything in the folder.
func verifyEmpty(kanaSettings *settings.Settings, consoleOutput *console.Console) error {
	if kanaSettings.Get("Type") == site.DefaultType {
		isEmpty, err := helpers.IsEmpty(kanaSettings.Get("Working"))
		if err != nil {
			return err
		}

		if !isEmpty && kanaSettings.GetBool("IsNew") {
			confirm := consoleOutput.PromptConfirm(
				"The current directory is not empty. Are you sure you want to try to install WordPress in this folder? This may cause the WordPress installation to fail.", //nolint: lll
				false)
			if !confirm {
				return fmt.Errorf("start aborted by user confirmation")
			}
		}
	}

	return nil
}
