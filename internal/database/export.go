package database

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana-cli/internal/site"
)

func Export(site *site.Site, args []string) (string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	exportFileName := fmt.Sprintf("kana-%s.sql", site.StaticConfig.SiteName)
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

	_, err = site.RunWPCli(exportCommand)
	if err != nil {
		return "", err
	}

	err = copyFile(path.Join(site.StaticConfig.SiteDirectory, "export.sql"), exportFile)
	if err != nil {
		return "", err
	}

	return exportFile, nil
}
