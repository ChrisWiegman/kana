package database

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana-cli/internal/app"
	"github.com/ChrisWiegman/kana-cli/internal/config"
)

func Export(kanaConfig *config.Config, args []string) (string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	exportFileName := fmt.Sprintf("kana-%s.sql", kanaConfig.Site.SiteName)
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

	site, err := app.NewSite(kanaConfig)
	if err != nil {
		return "", err
	}

	_, err = site.RunWPCli(exportCommand)
	if err != nil {
		return "", err
	}

	err = copyFile(path.Join(kanaConfig.Directories.Site, "export.sql"), exportFile)
	if err != nil {
		return "", err
	}

	return exportFile, nil
}
