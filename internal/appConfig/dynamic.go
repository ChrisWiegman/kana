package appConfig

import (
	"fmt"
	"os"
	"path"

	"github.com/aquasecurity/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ValidPHPVersions = []string{
	"site",
	"plugin",
	"theme",
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
	dynamicConfig.SetDefault("adminUser", "admin")
	dynamicConfig.SetDefault("adminPassword", "password")
	dynamicConfig.SetDefault("adminEmail", "admin@mykanasite.localhost")

	dynamicConfig.SetConfigName("kana")
	dynamicConfig.SetConfigType("json")
	dynamicConfig.AddConfigPath(path.Join(staticConfig.AppDirectory, "config"))

	err := dynamicConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			err = dynamicConfig.SafeWriteConfig()
			if err != nil {
				fmt.Println("error 1")
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

	t.AddRow("adminEmail", dynamicConfig.GetString("adminEmail"))
	t.AddRow("adminPassword", dynamicConfig.GetString("adminPassword"))
	t.AddRow("adminUser", dynamicConfig.GetString("adminUser"))
	t.AddRow("local", dynamicConfig.GetString("local"))
	t.AddRow("php", dynamicConfig.GetString("php"))
	t.AddRow("type", dynamicConfig.GetString("type"))
	t.AddRow("xdebug", dynamicConfig.GetString("xdebug"))

	t.Render()
}

func SetDynamicContent(md *cobra.Command, args []string, dynamicConfig *viper.Viper) {

}
