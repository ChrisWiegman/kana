package site

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ChrisWiegman/kana-dev/internal/console"
	"github.com/ChrisWiegman/kana-dev/internal/docker"
	"github.com/ChrisWiegman/kana-dev/internal/helpers"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
)

func (s *Site) ExportDatabase(args []string, consoleOutput *console.Console) (string, error) {
	isUsingSQLite, err := s.isUsingSQLite()
	if err != nil {
		return "", err
	}

	if isUsingSQLite {
		return "", fmt.Errorf("SQLite databases cannot be exported")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	exportFileName := fmt.Sprintf("kana-%s.sql", s.settings.Get("name"))
	exportFile := filepath.Join(cwd, exportFileName)

	if len(args) == 1 {
		exportFile = filepath.Join(cwd, args[0])
	}

	exportCommand := []string{
		"db",
		"export",
		"--add-drop-table",
		"/Site/export.sql",
	}

	code, output, err := s.WPCli(exportCommand, false, consoleOutput)
	if err != nil || code != 0 {
		errorMessage := ""

		if err != nil {
			errorMessage = err.Error()
		}

		return "", fmt.Errorf("database export failed: %s\n%s", errorMessage, output)
	}

	err = copyFile(filepath.Join(s.settings.Get("siteDirectory"), "export.sql"), exportFile)
	if err != nil {
		return "", err
	}

	return exportFile, nil
}

func (s *Site) ImportDatabase(file string, preserve bool, replaceDomain string, consoleOutput *console.Console) error {
	isUsingSQLite, err := s.isUsingSQLite()
	if err != nil {
		return err
	}

	if isUsingSQLite {
		return fmt.Errorf("SQLite databases cannot be imported")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	rawImportFile := filepath.Join(cwd, file)
	if _, err = os.Stat(rawImportFile); os.IsNotExist(err) {
		return fmt.Errorf("the specified sql file does not exist. Please enter a valid file to import")
	}

	kanaImportFile := filepath.Join(s.settings.Get("siteDirectory"), "import.sql")

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

		code, output, err = s.WPCli(dropCommand, false, consoleOutput)
		if err != nil || code != 0 {
			return fmt.Errorf("drop database failed: %s\n%s", err.Error(), output)
		}

		code, output, err = s.WPCli(createCommand, false, consoleOutput)
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

	code, output, err := s.WPCli(importCommand, false, consoleOutput)
	if err != nil || code != 0 {
		return fmt.Errorf("database import failed: %s\n%s", err.Error(), output)
	}

	if replaceDomain != "" {
		consoleOutput.Println("Replacing the old domain name")

		replaceCommand := []string{
			"search-replace",
			replaceDomain,
			s.settings.GetDomain(),
			"--all-tables",
		}

		code, output, err := s.WPCli(replaceCommand, false, consoleOutput)
		if err != nil || code != 0 {
			return fmt.Errorf("replace domain failed failed: %s\n%s", err.Error(), output)
		}
	}

	return nil
}

func (s *Site) getDatabaseContainer(databaseDir string, appContainers []docker.ContainerConfig) []docker.ContainerConfig {
	isUsingSQLite, err := s.isUsingSQLite()
	if err != nil {
		return appContainers
	}

	if isUsingSQLite {
		return appContainers
	}

	envVars := []string{
		"MARIADB_ROOT_PASSWORD=password",
		"MARIADB_DATABASE=wordpress",
		"MARIADB_USER=wordpress",
		"MARIADB_PASSWORD=wordpress",
	}

	if s.settings.Get("database") == "mysql" {
		envVars = []string{
			"MYSQL_ROOT_PASSWORD=password",
			"MYSQL_DATABASE=wordpress",
			"MYSQL_USER=wordpress",
			"MYSQL_PASSWORD=wordpress",
		}
	}

	databaseContainer := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana-%s-database", s.settings.Get("name")),
		Image:       fmt.Sprintf("%s:%s", s.settings.Get("database"), s.settings.Get("databaseVersion")),
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana-%s-database", s.settings.Get("name")),
		Ports: []docker.ExposedPorts{
			{Port: "3306", Protocol: "tcp"},
		},
		Env: envVars,
		Labels: map[string]string{
			"kana.type": "database",
			"kana.site": s.settings.Get("name"),
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

func (s *Site) getDatabaseDirectory() (databaseDirectory string, err error) {
	databaseDirectory = filepath.Join(s.settings.Get("siteDirectory"), "database")

	err = os.MkdirAll(databaseDirectory, os.FileMode(defaultDirPermissions))
	if err != nil {
		return "", err
	}

	return databaseDirectory, err
}

// getDatabasePort returns the public port for the database attached to the current site.
func (s *Site) getDatabasePort() string {
	containers, _ := s.dockerClient.ContainerList(s.settings.Get("name"))
	var databasePort types.Port

	for i := range containers {
		if containers[i].Image == fmt.Sprintf("%s:%s", s.settings.Get("database"), s.settings.Get("databaseVersion")) {
			databasePort = containers[i].Ports[0]
		}
	}

	return strconv.Itoa(int(databasePort.PublicPort))
}

func (s *Site) maybeSetupSQLite() error {
	isUsingSQLite, err := s.isUsingSQLite()
	if err != nil {
		return err
	}

	if !isUsingSQLite {
		return nil
	}

	file, err := helpers.DownloadFile(
		"https://downloads.wordpress.org/plugin/sqlite-database-integration.zip",
		s.settings.Get("workingDirectory"))
	if err != nil {
		return err
	}

	err = helpers.UnZipFile(
		filepath.Join(s.settings.Get("workingDirectory"), file),
		filepath.Join(s.settings.Get("workingDirectory"), "wp-content", "plugins"))
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join(s.settings.Get("workingDirectory"), file))
	if err != nil {
		return err
	}

	return helpers.CopyFile(
		filepath.Join(
			s.settings.Get("workingDirectory"), "wp-content", "plugins", "sqlite-database-integration", "db.copy"),
		filepath.Join(s.settings.Get("workingDirectory"), "wp-content", "db.php"))
}

func (s *Site) isUsingSQLite() (bool, error) {
	output, err := s.WordPress("echo $KANA_SQLITE", false, false)
	if err != nil {
		return false, err
	}

	if strings.Contains(output.StdOut, "true") {
		return true, nil
	}

	return false, nil
}

// verifySite verifies if a site is up and running without error.
func (s *Site) verifyDatabase(consoleOutput *console.Console) error {
	checkCommand := []string{
		"db",
		"check",
	}

	databaseOK := false
	checkAttempt := 0

	for !databaseOK {
		code, _, err := s.WPCli(checkCommand, false, consoleOutput)
		if err != nil || code != 0 {
			checkAttempt++ // Increment the check attempt counter
			time.Sleep(time.Second)
		} else {
			return nil
		}

		if checkAttempt == s.maxVerificationRetries {
			return fmt.Errorf("database verification failed")
		}
	}

	return nil
}
