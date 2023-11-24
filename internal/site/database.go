package site

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/docker"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
)

func (s *Site) ExportDatabase(args []string, consoleOutput *console.Console) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	exportFileName := fmt.Sprintf("kana-%s.sql", s.Settings.Name)
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

	code, output, err := s.RunWPCli(exportCommand, consoleOutput)
	if err != nil || code != 0 {
		return "", fmt.Errorf("database export failed: %s\n%s", err.Error(), output)
	}

	err = copyFile(path.Join(s.Settings.SiteDirectory, "export.sql"), exportFile)
	if err != nil {
		return "", err
	}

	return exportFile, nil
}

func (s *Site) ImportDatabase(file string, preserve bool, replaceDomain string, consoleOutput *console.Console) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	rawImportFile := path.Join(cwd, file)
	if _, err = os.Stat(rawImportFile); os.IsNotExist(err) {
		return fmt.Errorf("the specified sql file does not exist. Please enter a valid file to import")
	}

	kanaImportFile := path.Join(s.Settings.SiteDirectory, "import.sql")

	err = copyFile(rawImportFile, kanaImportFile)
	if err != nil {
		return err
	}

	if !preserve {
		consoleOutput.Println("Dropping the existing database.")

		dropCommand := []string{
			"db",
			"drop",
			"--yes",
		}

		createCommand := []string{
			"db",
			"create",
		}

		var code int64
		var output string

		code, output, err = s.RunWPCli(dropCommand, consoleOutput)
		if err != nil || code != 0 {
			return fmt.Errorf("drop database failed: %s\n%s", err.Error(), output)
		}

		code, output, err = s.RunWPCli(createCommand, consoleOutput)
		if err != nil || code != 0 {
			return fmt.Errorf("create database failed: %s\n%s", err.Error(), output)
		}
	}

	consoleOutput.Println("Importing the database file.")

	importCommand := []string{
		"db",
		"import",
		"/Site/import.sql",
	}

	code, output, err := s.RunWPCli(importCommand, consoleOutput)
	if err != nil || code != 0 {
		return fmt.Errorf("database import failed: %s\n%s", err.Error(), output)
	}

	if replaceDomain != "" {
		consoleOutput.Println("Replacing the old domain name")

		replaceCommand := []string{
			"search-replace",
			replaceDomain,
			s.Settings.SiteDomain,
			"--all-tables",
		}

		code, output, err := s.RunWPCli(replaceCommand, consoleOutput)
		if err != nil || code != 0 {
			return fmt.Errorf("replace domain failed failed: %s\n%s", err.Error(), output)
		}
	}

	return nil
}

func (s *Site) getDatabaseContainer(databaseDir string, appContainers []docker.ContainerConfig) []docker.ContainerConfig {
	databaseContainer := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana-%s-database", s.Settings.Name),
		Image:       fmt.Sprintf("mariadb:%s", s.Settings.MariaDB),
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana-%s-database", s.Settings.Name),
		Ports: []docker.ExposedPorts{
			{Port: "3306", Protocol: "tcp"},
		},
		Env: []string{
			"MARIADB_ROOT_PASSWORD=password",
			"MARIADB_DATABASE=wordpress",
			"MARIADB_USER=wordpress",
			"MARIADB_PASSWORD=wordpress",
		},
		Labels: map[string]string{
			"kana.type": "database",
			"kana.site": s.Settings.Name,
		},
		Volumes: []mount.Mount{
			{ // Maps a database folder to the MySQL container for persistence
				Type:   mount.TypeBind,
				Source: databaseDir,
				Target: "/var/lib/mysql",
			},
		},
	}

	appContainers = append(appContainers, databaseContainer)

	return appContainers
}

// getDatabasePort returns the public port for the database attached to the current site.
func (s *Site) getDatabasePort() string {
	containers, _ := s.dockerClient.ContainerList(s.Settings.Name)
	var databasePort types.Port

	for i := range containers {
		if containers[i].Image == fmt.Sprintf("mariadb:%s", s.Settings.MariaDB) {
			databasePort = containers[i].Ports[0]
		}
	}

	return strconv.Itoa(int(databasePort.PublicPort))
}
