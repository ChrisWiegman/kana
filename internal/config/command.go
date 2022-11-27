package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aquasecurity/table"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

func (c *Config) GetDynamicContentItem(md *cobra.Command, args []string) (string, error) {

	if !c.Global.Viper.IsSet(args[0]) {
		return "", fmt.Errorf("invalid setting. Please enter a valid key to get")
	}

	return c.Global.Viper.GetString(args[0]), nil
}

func (c *Config) ListDynamicContent() {

	t := table.New(os.Stdout)

	t.SetHeaders("Key", "Value")

	t.AddRow("admin.email", c.Global.Viper.GetString("admin.email"))
	t.AddRow("admin.password", c.Global.Viper.GetString("admin.password"))
	t.AddRow("admnin.username", c.Global.Viper.GetString("admin.username"))
	t.AddRow("local", c.Global.Viper.GetString("local"))
	t.AddRow("php", c.Global.Viper.GetString("php"))
	t.AddRow("type", c.Global.Viper.GetString("type"))
	t.AddRow("xdebug", c.Global.Viper.GetString("xdebug"))

	t.Render()
}

func (c *Config) SetDynamicContent(md *cobra.Command, args []string) error {

	if !c.Global.Viper.IsSet(args[0]) {
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
		c.Global.Viper.Set(args[0], boolVal)
		return c.Global.Viper.WriteConfig()
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

	c.Global.Viper.Set(args[0], args[1])

	return c.Global.Viper.WriteConfig()
}
