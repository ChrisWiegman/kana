package appConfig

import (
	"fmt"
	"path"

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
