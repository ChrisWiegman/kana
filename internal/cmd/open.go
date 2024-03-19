package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var openDatabaseFlag, openMailpitFlag, openSiteFlag, openAdminFlag bool

func open(consoleOutput *console.Console, kanaSite *site.Site) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open the current site in your browser.",
		Run: func(cmd *cobra.Command, args []string) {
			err := kanaSite.EnsureDocker(consoleOutput)
			if err != nil {
				consoleOutput.Error(err)
			}

			// Default to opening the site if no flags are specified
			if !cmd.Flags().Lookup("database").Changed &&
				!cmd.Flags().Lookup("mailpit").Changed &&
				!cmd.Flags().Lookup("site").Changed &&
				!cmd.Flags().Lookup("admin").Changed {
				openSiteFlag = true
			}

			// Open the site in the user's default browser,
			err = kanaSite.OpenSite(openDatabaseFlag, openMailpitFlag, openSiteFlag, openAdminFlag, consoleOutput)
			if err != nil {
				consoleOutput.Error(fmt.Errorf("an error occurred and we can't open the requested resource: %s", err))
			}

			consoleOutput.Success(
				fmt.Sprintf(
					"Your site, %s, has been opened in your default browser.",
					consoleOutput.Bold(
						consoleOutput.Blue(
							kanaSite.Settings.Name))))
		},
		Args: cobra.NoArgs,
	}

	commandsRequiringSite = append(commandsRequiringSite, cmd.Use)
	cmd.Flags().BoolVarP(
		&openDatabaseFlag,
		"database",
		"d",
		false,
		"Opens the database for the current or specified Kana site with phpMyAdmin in your default browser")
	cmd.Flags().BoolVarP(
		&openMailpitFlag,
		"mailpit",
		"m",
		false,
		"Opens the Mailpit UI for the current or specified Kana site in your default browser")
	cmd.Flags().BoolVarP(&openSiteFlag, "site", "s", false, "Opens the current or specified Kana site in your default browser")
	cmd.Flags().BoolVarP(
		&openAdminFlag,
		"admin",
		"a",
		false,
		"Opens the current or specified Kana site's WordPress dashboard in your default browser")

	cmd.Flags().SetNormalizeFunc(aliasPhpMyAdminFlag)

	return cmd
}

func aliasPhpMyAdminFlag(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if name == "phpmyadmin" {
		name = "database"
	}

	return pflag.NormalizedName(name)
}
