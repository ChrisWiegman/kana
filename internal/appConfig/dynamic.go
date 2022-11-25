package appConfig

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/aquasecurity/table"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ValidPHPVersions = []string{
	"7.4",
	"8.0",
	"8.1",
}

var ValidTypes = []string{
	"site",
	"plugin",
	"theme",
}

func GetDynamicContent(staticConfig StaticConfig) (*viper.Viper, error) {

	dynamicConfig := viper.New()

	dynamicConfig.SetDefault("xdebug", false)
	dynamicConfig.SetDefault("type", "site")
	dynamicConfig.SetDefault("local", false)
	dynamicConfig.SetDefault("php", "7.4")
	dynamicConfig.SetDefault("admin.username", "admin")
	dynamicConfig.SetDefault("admin.password", "password")
	dynamicConfig.SetDefault("admin.email", "admin@mykanasite.localhost")

	dynamicConfig.SetConfigName("kana")
	dynamicConfig.SetConfigType("json")
	dynamicConfig.AddConfigPath(path.Join(staticConfig.AppDirectory, "config"))

	err := dynamicConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			err = dynamicConfig.SafeWriteConfig()
			if err != nil {
				return dynamicConfig, err
			}
		} else {
			return dynamicConfig, err
		}
	}

	changeConfig := false

	// Reset default "site" type if there's an invalid type in the config file
	if !CheckString(dynamicConfig.GetString("type"), ValidTypes) {
		changeConfig = true
		dynamicConfig.Set("type", "site")
	}

	// Reset default php version if there's an invalid version in the config file
	if !CheckString(dynamicConfig.GetString("php"), ValidPHPVersions) {
		changeConfig = true
		dynamicConfig.Set("php", "7.4")
	}

	if changeConfig {
		err = dynamicConfig.WriteConfig()
		if err != nil {
			return dynamicConfig, err
		}
	}

	return dynamicConfig, nil
}

func ListDynamicContent(dynamicConfig *viper.Viper) {

	t := table.New(os.Stdout)

	t.SetHeaders("Key", "Value")

	t.AddRow("admin.email", dynamicConfig.GetString("admin.email"))
	t.AddRow("admin.password", dynamicConfig.GetString("admin.password"))
	t.AddRow("admnin.username", dynamicConfig.GetString("admin.username"))
	t.AddRow("local", dynamicConfig.GetString("local"))
	t.AddRow("php", dynamicConfig.GetString("php"))
	t.AddRow("type", dynamicConfig.GetString("type"))
	t.AddRow("xdebug", dynamicConfig.GetString("xdebug"))

	t.Render()
}

func GetDynamicContentItem(md *cobra.Command, args []string, dynamicConfig *viper.Viper) (string, error) {

	if !dynamicConfig.IsSet(args[0]) {
		return "", fmt.Errorf("invalid setting. Please enter a valid key to get")
	}

	return dynamicConfig.GetString(args[0]), nil
}

func SetDynamicContent(md *cobra.Command, args []string, dynamicConfig *viper.Viper) error {

	if !dynamicConfig.IsSet(args[0]) {
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
		dynamicConfig.Set(args[0], boolVal)
		return dynamicConfig.WriteConfig()
	case "php":
		if !CheckString(args[1], ValidPHPVersions) {
			err = fmt.Errorf("please choose a valid php version")
		}
	case "type":
		if !CheckString(args[1], ValidTypes) {
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

	dynamicConfig.Set(args[0], args[1])

	return dynamicConfig.WriteConfig()
}
