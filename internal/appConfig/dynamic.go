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

type DynamicConfig struct {
	SiteXdebug    bool
	SiteLocal     bool
	SiteType      string
	PHPVersion    string
	AdminUsername string
	AdminPassword string
	AdminEmail    string
}

func GetDynamicContent(staticConfig StaticConfig) (DynamicConfig, error) {

	viperConfig, err := loadDynamicConfig(staticConfig.AppDirectory)
	if err != nil {
		return DynamicConfig{}, err
	}

	dynamicConfig := DynamicConfig{
		SiteXdebug:    viperConfig.GetBool("xdebug"),
		SiteLocal:     viperConfig.GetBool("xdebug"),
		SiteType:      viperConfig.GetString("type"),
		PHPVersion:    viperConfig.GetString("php"),
		AdminUsername: viperConfig.GetString("adminUser"),
		AdminPassword: viperConfig.GetString("adminPassword"),
		AdminEmail:    viperConfig.GetString("adminEmail"),
	}

	return dynamicConfig, nil

}

func loadDynamicConfig(appDirectory string) (*viper.Viper, error) {

	defaultConfig := viper.New()

	defaultConfig.SetDefault("xdebug", false)
	defaultConfig.SetDefault("type", "site")
	defaultConfig.SetDefault("local", false)
	defaultConfig.SetDefault("php", "7.4")
	defaultConfig.SetDefault("adminUser", "admin")
	defaultConfig.SetDefault("adminPassword", "password")
	defaultConfig.SetDefault("adminEmail", "admin@mykanasite.localhost")

	defaultConfig.SetConfigName("kana")
	defaultConfig.SetConfigType("json")
	defaultConfig.AddConfigPath(path.Join(appDirectory, "config"))

	err := defaultConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			err = defaultConfig.SafeWriteConfig()
			if err != nil {
				fmt.Println("error 1")
				return defaultConfig, err
			}
		} else {
			return defaultConfig, err
		}
	}

	changeConfig := false

	// Reset default "site" type if there's an invalid type in the config file
	if !CheckString(defaultConfig.GetString("type"), ValidTypes) {
		changeConfig = true
		defaultConfig.Set("type", "site")
	}

	// Reset default php version if there's an invalid version in the config file
	if !CheckString(defaultConfig.GetString("php"), ValidPHPVersions) {
		changeConfig = true
		defaultConfig.Set("php", "7.4")
	}

	if changeConfig {
		err = defaultConfig.WriteConfig()
		if err != nil {
			return defaultConfig, err
		}
	}

	return defaultConfig, nil

}
