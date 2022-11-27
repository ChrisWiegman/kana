package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aquasecurity/table"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

func (c *Config) GetGlobalSetting(md *cobra.Command, args []string) (string, error) {

	if !c.Global.viper.IsSet(args[0]) {
		return "", fmt.Errorf("invalid setting. Please enter a valid key to get")
	}

	return c.Global.viper.GetString(args[0]), nil
}

func (c *Config) ListConfig() {

	t := table.New(os.Stdout)

	t.SetHeaders("Key", "Global Value", "Local Value")

	t.AddRow("admin.email", c.Global.AdminEmail)
	t.AddRow("admin.password", c.Global.AdminPassword)
	t.AddRow("admnin.username", c.Global.AdminUsername)
	t.AddRow("local", strconv.FormatBool(c.Global.Local), strconv.FormatBool(c.Local.Local))
	t.AddRow("php", c.Global.PHP, c.Local.PHP)
	t.AddRow("type", c.Global.Type, c.Local.Type)
	t.AddRow("xdebug", strconv.FormatBool(c.Global.Xdebug), strconv.FormatBool(c.Local.Local))

	plugins := strings.Join(c.Local.Plugins, "\n")

	t.AddRow("plugins", "", plugins)

	t.Render()
}

func (c *Config) SetGlobalSetting(md *cobra.Command, args []string) error {

	if !c.Global.viper.IsSet(args[0]) {
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
		boolVal, err := strconv.ParseBool(args[1])
		if err != nil {
			return err
		}
		c.Global.viper.Set(args[0], boolVal)
		return c.Global.viper.WriteConfig()
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

	c.Global.viper.Set(args[0], args[1])

	return c.Global.viper.WriteConfig()
}
