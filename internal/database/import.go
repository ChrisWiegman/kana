package database

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana-cli/internal/app"
	"github.com/ChrisWiegman/kana-cli/internal/config"
	"github.com/ChrisWiegman/kana-cli/internal/console"
)

func Import(kanaConfig *config.Config, file string, preserve bool, replaceDomain string) error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	rawImportFile := path.Join(cwd, file)
	if _, err = os.Stat(rawImportFile); os.IsNotExist(err) {
		return fmt.Errorf("the specified sql file does not exist. Please enter a valid file to import")
	}

	kanaImportFile := path.Join(kanaConfig.Directories.Site, "import.sql")

	err = copyFile(rawImportFile, kanaImportFile)
	if err != nil {
		return err
	}

	site, err := app.NewSite(kanaConfig)
	if err != nil {
		return err
	}

	if !preserve {

		console.Println("Dropping the existing database.")

		dropCommand := []string{
			"db",
			"drop",
			"--yes",
		}

		createCommand := []string{
			"db",
			"create",
		}

		_, err = site.RunWPCli(dropCommand)
		if err != nil {
			return err
		}

		_, err = site.RunWPCli(createCommand)
		if err != nil {
			return err
		}
	}

	console.Println("Importing the database file.")

	importCommand := []string{
		"db",
		"import",
		"/Site/import.sql",
	}

	_, err = site.RunWPCli(importCommand)
	if err != nil {
		return err
	}

	if replaceDomain != "" {

		console.Println("Replacing the old domain name")

		replaceCommand := []string{
			"search-replace",
			replaceDomain,
			fmt.Sprintf("%s.%s", kanaConfig.Site.SiteName, kanaConfig.App.AppDomain),
			"--all-tables",
		}

		_, err := site.RunWPCli(replaceCommand)
		if err != nil {
			return err
		}
	}

	return nil
}
