package config

import (
	"path"

	"github.com/spf13/viper"
)

var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"
var appDomain = "sites.kana.li"
var configFolderName = ".config/kana"

type AppConfig struct {
	Xdebug        bool
	Type          string
	Local         bool
	PHP           string
	AdminUsername string
	AdminPassword string
	AdminEmail    string
	AppDomain     string
	AppDirectory  string
	RootKey       string
	RootCert      string
	SiteCert      string
	SiteKey       string
	Viper         *viper.Viper
}

// getAppConfig gets config information that transcends sites such as app and default settings
func (c *Config) getAppConfig() (AppConfig, error) {

	appConfig := AppConfig{
		AppDomain: appDomain,
		RootKey:   rootKey,
		RootCert:  rootCert,
		SiteCert:  siteCert,
		SiteKey:   siteKey,
	}

	appDirectory, err := getAppDirectory()
	if err != nil {
		return appConfig, err
	}

	appConfig.AppDirectory = appDirectory

	dynamicConfig, err := c.loadAppConfig(appDirectory)
	if err != nil {
		return appConfig, err
	}

	appConfig.Viper = dynamicConfig
	appConfig.Xdebug = dynamicConfig.GetBool("xdebug")
	appConfig.Local = dynamicConfig.GetBool("local")
	appConfig.AdminEmail = dynamicConfig.GetString("admin.email")
	appConfig.AdminPassword = dynamicConfig.GetString("admin.password")
	appConfig.AdminUsername = dynamicConfig.GetString("admin.username")
	appConfig.PHP = dynamicConfig.GetString("php")
	appConfig.Type = dynamicConfig.GetString("type")

	return appConfig, err

}

// loadAppConfig loads the app config using viper and sets defaults
func (c *Config) loadAppConfig(appDirectory string) (*viper.Viper, error) {

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
	dynamicConfig.AddConfigPath(path.Join(appDirectory, "config"))

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
	if !CheckString(dynamicConfig.GetString("type"), validTypes) {
		changeConfig = true
		dynamicConfig.Set("type", "site")
	}

	// Reset default php version if there's an invalid version in the config file
	if !CheckString(dynamicConfig.GetString("php"), validPHPVersions) {
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
