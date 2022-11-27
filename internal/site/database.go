package site

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana-cli/pkg/console"
)

func (s *Site) ExportDatabase(args []string) (string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	exportFileName := fmt.Sprintf("kana-%s.sql", s.Config.Name)
	exportFile := path.Join(cwd, exportFileName)

	if len(args) == 1 {
		exportFile = path.Join(cwd, args[0])
	}

	exportCommand := []string{
		"db",
		"export",
		"--add-drop-table",
		"/Site/export.sql",
	}

	_, err = s.RunWPCli(exportCommand)
	if err != nil {
		return "", err
	}

	err = copyFile(path.Join(s.Config.SiteDirectory, "export.sql"), exportFile)
	if err != nil {
		return "", err
	}

	return exportFile, nil
}

func (s *Site) ImportDatabase(file string, preserve bool, replaceDomain string) error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	rawImportFile := path.Join(cwd, file)
	if _, err = os.Stat(rawImportFile); os.IsNotExist(err) {
		return fmt.Errorf("the specified sql file does not exist. Please enter a valid file to import")
	}

	kanaImportFile := path.Join(s.Config.SiteDirectory, "import.sql")

	err = copyFile(rawImportFile, kanaImportFile)
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

		_, err = s.RunWPCli(dropCommand)
		if err != nil {
			return err
		}

		_, err = s.RunWPCli(createCommand)
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

	_, err = s.RunWPCli(importCommand)
	if err != nil {
		return err
	}

	if replaceDomain != "" {

		console.Println("Replacing the old domain name")

		replaceCommand := []string{
			"search-replace",
			replaceDomain,
			s.Config.SiteDomain,
			"--all-tables",
		}

		_, err := s.RunWPCli(replaceCommand)
		if err != nil {
			return err
		}
	}

	return nil
}
