package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/aquasecurity/table"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

// GetGlobalSetting Retrieves a global setting for the "config" command
func (s *Settings) GetGlobalSetting(md *cobra.Command, args []string) (string, error) {
	if !s.global.IsSet(args[0]) {
		return "", fmt.Errorf("invalid setting. Please enter a valid key to get")
	}

	return s.global.GetString(args[0]), nil
}

// ListSettings Lists all settings for the config command
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
	t.AddRow("local", consoleOutput.Bold(s.global.GetString("local")), consoleOutput.Bold(s.local.GetString("local")))
	t.AddRow("mailpit", consoleOutput.Bold(s.global.GetString("mailpit")), consoleOutput.Bold(s.local.GetString("mailpit")))
	t.AddRow("php", consoleOutput.Bold(s.global.GetString("php")), consoleOutput.Bold(s.local.GetString("php")))
	t.AddRow("phpmyadmin", consoleOutput.Bold(s.global.GetString("phpmyadmin")), consoleOutput.Bold(s.local.GetString("phpmyadmin")))

	boldPlugins := []string{}

	for _, plugin := range s.Plugins {
		boldPlugins = append(boldPlugins, consoleOutput.Bold(plugin))
	}

	plugins := consoleOutput.Bold(strings.Join(boldPlugins, "\n"))

	t.AddRow("plugins", "", plugins)

	t.AddRow("ssl", consoleOutput.Bold(s.global.GetString("ssl")), consoleOutput.Bold(s.local.GetString("ssl")))
	t.AddRow("type", consoleOutput.Bold(s.global.GetString("type")), consoleOutput.Bold(s.local.GetString("type")))
	t.AddRow("xdebug", consoleOutput.Bold(s.global.GetString("xdebug")), consoleOutput.Bold(s.local.GetString("xdebug")))

	t.Render()
}

// printJSONSettings Prints out all settings in JSON format
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

// SetGlobalSetting Sets a global setting for the "config" command
func (s *Settings) SetGlobalSetting(md *cobra.Command, args []string) error {
	if !s.global.IsSet(args[0]) {
		return fmt.Errorf("invalid setting. Please enter a valid key to set")
	}

	validate := validator.New()
	var err error

	switch args[0] {
	case "local", "xdebug", "ssl":
		err = validate.Var(args[1], "boolean")
		if err != nil {
			return err
		}

		var boolVal bool

		boolVal, err = strconv.ParseBool(args[1])
		if err != nil {
			return err
		}
		s.global.Set(args[0], boolVal)
		return s.global.WriteConfig()
	case "php":
		if !isValidString(args[1], validPHPVersions) {
			err = fmt.Errorf("please choose a valid php version")
		}
	case "type":
		if !isValidString(args[1], validTypes) {
			err = fmt.Errorf("please choose a valid project type")
		}
	case "admin.email":
		err = validate.Var(args[1], "email")
	case "admin.password":
		err = validate.Var(args[1], "alphanumunicode")
	case "admin.username":
		err = validate.Var(args[1], "alpha")
	default:
		err = validate.Var(args[1], "boolean")
	}

	if err != nil {
		return err
	}

	s.global.Set(args[0], args[1])

	return s.global.WriteConfig()
}
