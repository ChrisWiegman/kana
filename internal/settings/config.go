package settings

import (
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
func (s *Settings) ListSettings() {
	t := table.New(os.Stdout)

	t.SetHeaders("Setting", "Global Value", "Local Value")

	t.AddRow("admin.email", console.Bold(s.global.GetString("admin.email")))
	t.AddRow("admin.password", console.Bold(s.global.GetString("admin.password")))
	t.AddRow("admnin.username", console.Bold(s.global.GetString("admin.username")))
	t.AddRow("dockerSockFile", console.Bold(s.global.GetString("dockerSockFile")))
	t.AddRow("local", console.Bold(s.global.GetString("local")), console.Bold(s.local.GetString("local")))
	t.AddRow("php", console.Bold(s.global.GetString("php")), console.Bold(s.local.GetString("php")))
	t.AddRow("phpmyadmin", console.Bold(s.global.GetString("phpmyadmin")), console.Bold(s.local.GetString("phpmyadmin")))
	t.AddRow("type", console.Bold(s.global.GetString("type")), console.Bold(s.local.GetString("type")))
	t.AddRow("xdebug", console.Bold(s.global.GetString("xdebug")), console.Bold(s.local.GetString("xdebug")))

	boldPlugins := []string{}

	for _, plugin := range s.Plugins {
		boldPlugins = append(boldPlugins, console.Bold(plugin))
	}

	plugins := console.Bold(strings.Join(boldPlugins, "\n"))

	t.AddRow("plugins", "", plugins)

	t.Render()
}

// SetGlobalSetting Sets a global setting for the "config" command
func (s *Settings) SetGlobalSetting(md *cobra.Command, args []string) error {
	if !s.global.IsSet(args[0]) {
		return fmt.Errorf("invalid setting. Please enter a valid key to set")
	}

	validate := validator.New()
	var err error

	switch args[0] {
	case "local", "xdebug":
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
