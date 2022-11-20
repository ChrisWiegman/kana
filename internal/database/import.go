package database

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/ChrisWiegman/kana-cli/internal/site"
)

func Import(site *site.Site, file string, preserve bool, replaceDomain string) error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	importFile := path.Join(cwd, file)
	if _, err = os.Stat(importFile); os.IsNotExist(err) {
		return fmt.Errorf("the specified sql file does not exist. Please enter a valid file to import")
	}

	source, err := os.Open(importFile)

	if err != nil {
		return err
	}
	defer source.Close()

	kanaImportFile := path.Join(site.StaticConfig.SiteDirectory, "import.sql")

	destination, err := os.Create(kanaImportFile)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	if !preserve {

		dropCommand := []string{
			"db",
			"drop",
			"--yes",
		}

		createCommand := []string{
			"db",
			"create",
		}

		out, err := site.RunWPCli(dropCommand)
		if err != nil {
			return err
		}

		fmt.Println(out)

		out, err = site.RunWPCli(createCommand)
		if err != nil {
			return err
		}

		fmt.Println(out)
	}

	importCommand := []string{
		"db",
		"import",
		"/Site/import.sql",
	}

	out, err := site.RunWPCli(importCommand)
	if err != nil {
		return err
	}

	fmt.Println(out)

	if replaceDomain != "" {

		replaceCommand := []string{
			"search-replace",
			replaceDomain,
			site.StaticConfig.AppDomain,
			"--all-tables",
		}

		out, err := site.RunWPCli(replaceCommand)
		if err != nil {
			return err
		}

		fmt.Println(out)
	}

	return nil
}
