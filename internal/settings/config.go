package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/helpers"

	"github.com/aquasecurity/table"
	"github.com/go-playground/validator/v10"
)

// GetGlobalSetting Retrieves a global setting for the "config" command.
func (s *Settings) GetGlobalSetting(args []string) (string, error) {
	if !s.global.IsSet(args[0]) {
		return "", fmt.Errorf("invalid setting. Please enter a valid key to get")
	}

	return s.global.GetString(args[0]), nil
}

// ListSettings Lists all settings for the config command.
func (s *Settings) ListSettings(consoleOutput *console.Console) {
	if consoleOutput.JSON {
		s.printJSONSettings()
		return
	}

	t := table.New(os.Stdout)

	t.SetHeaders("Setting", "Global Value", "Local Value")

	t.AddRow("admin.email", consoleOutput.Bold(s.global.GetString("admin.email")))
	t.AddRow("admin.password", consoleOutput.Bold(s.global.GetString("admin.password")))
	t.AddRow("admnin.username", consoleOutput.Bold(s.global.GetString("admin.username")))
	t.AddRow(
		"imageUpdateDays",
		consoleOutput.Bold(s.global.GetString("imageUpdateDays")),
		consoleOutput.Bold(s.local.GetString("imageUpdateDays")))
	t.AddRow("multisite", consoleOutput.Bold(s.global.GetString("multisite")), consoleOutput.Bold(s.local.GetString("multisite")))
	t.AddRow("mailpit", consoleOutput.Bold(s.global.GetString("mailpit")), consoleOutput.Bold(s.local.GetString("mailpit")))
	t.AddRow("php", consoleOutput.Bold(s.global.GetString("php")), consoleOutput.Bold(s.local.GetString("php")))
	t.AddRow("mariadb", consoleOutput.Bold(s.global.GetString("mariadb")), consoleOutput.Bold(s.local.GetString("mariadb")))
	t.AddRow(
		"databaseClient",
		consoleOutput.Bold(s.global.GetString("databaseClient")),
		consoleOutput.Bold(s.local.GetString("databaseClient")))

	boldPlugins := []string{}

	for _, plugin := range s.Plugins {
		boldPlugins = append(boldPlugins, consoleOutput.Bold(plugin))
	}

	plugins := consoleOutput.Bold(strings.Join(boldPlugins, "\n"))

	t.AddRow("plugins", "", plugins)

	t.AddRow(
		"removeDefaultPlugins",
		consoleOutput.Bold(s.global.GetString("removeDefaultPlugins")),
		consoleOutput.Bold(s.local.GetString("removeDefaultPlugins")))
	t.AddRow("ssl", consoleOutput.Bold(s.global.GetString("ssl")), consoleOutput.Bold(s.local.GetString("ssl")))
	t.AddRow("type", consoleOutput.Bold(s.global.GetString("type")), consoleOutput.Bold(s.local.GetString("type")))
	t.AddRow("activate", consoleOutput.Bold(s.global.GetString("activate")), consoleOutput.Bold(s.local.GetString("activate")))
	t.AddRow("wpdebug", consoleOutput.Bold(s.global.GetString("wpdebug")), consoleOutput.Bold(s.local.GetString("wpdebug")))
	t.AddRow("xdebug", consoleOutput.Bold(s.global.GetString("xdebug")), consoleOutput.Bold(s.local.GetString("xdebug")))
	t.AddRow("scriptdebug", consoleOutput.Bold(s.global.GetString("scriptdebug")), consoleOutput.Bold(s.local.GetString("scriptdebug")))
	t.AddRow("environment", consoleOutput.Bold(s.global.GetString("environment")), consoleOutput.Bold(s.local.GetString("environment")))

	t.Render()
}

// SetGlobalSetting Sets a global setting for the "config" command.
func (s *Settings) SetGlobalSetting(args []string) error {
	if !s.global.IsSet(args[0]) {
		return fmt.Errorf("invalid setting. Please enter a valid key to set")
	}

	err := s.validateSetting(args[0], args[1])
	if err != nil {
		return err
	}

	s.global.Set(args[0], args[1])

	return s.global.WriteConfig()
}

// printJSONSettings Prints out all settings in JSON format.
func (s *Settings) printJSONSettings() {
	type JSONSettings struct {
		Global, Local map[string]interface{}
	}

	settings := JSONSettings{
		Global: s.global.AllSettings(),
		Local:  s.local.AllSettings(),
	}

	str, _ := json.Marshal(settings)

	fmt.Println(string(str))
}

// validateSetting validates the values in saved settings.
func (s *Settings) validateSetting(setting, value string) error { //nolint:gocyclo
	validate := validator.New()

	switch setting {
	case "php":
		if docker.ValidateImage("wordpress", fmt.Sprintf("php%s", value)) != nil {
			return fmt.Errorf(
				"the PHP version in your configuration, %s, is invalid. See https://hub.docker.com/_/wordpress for a list of supported versions",
				value)
		}
	case "mariadb":
		if docker.ValidateImage("mariadb", value) != nil {
			return fmt.Errorf(
				"the MariaDB version in your configuration, %s, is invalid. See https://hub.docker.com/_/mariadb for a list of supported versions",
				value)
		}
	case "type":
		if !helpers.IsValidString(value, validTypes) {
			return fmt.Errorf("the type you selected, %s, is not a valid type. You must use either `site`, `plugin` or `theme`", setting)
		}
	case "admin.email":
		return validate.Var(value, "email")
	case "admin.password":
		return validate.Var(value, "alphanumunicode")
	case "admin.username":
		return validate.Var(value, "alpha")
	case "imageUpdateDays":
	case "imageupdatedays":
		return validate.Var(value, "gte=0")
	case "databaseClient":
	case "databaseclient":
		if !helpers.IsValidString(value, validDatabaseClients) {
			return fmt.Errorf("the database client, %s, is not a valid client. You must use either `phpmyadmin` or `tableplus`", setting)
		}
	case "multisite":
		if !helpers.IsValidString(value, validMultisiteTypes) {
			return fmt.Errorf("the multisite type, %s, is not a valid type. You must use either `none`, `subdomain` or `subdirectory`", setting)
		}
	case "environment":
		if !helpers.IsValidString(value, validEnvironmentTypes) {
			return fmt.Errorf("the environment, %s, is not valid. You must use either `local`, `development`, `staging` or `production`", setting)
		}
	default:
		err := validate.Var(value, "boolean")
		if err != nil {
			return fmt.Errorf("the setting, %s, must be either true or false", setting)
		}
	}

	return nil
}
