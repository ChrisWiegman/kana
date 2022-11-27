package config

import (
	"path"

	"github.com/spf13/viper"
)

var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"
var domain = "sites.kana.li"
var configFolderName = ".config/kana"

type GlobalConfig struct {
	Xdebug        bool
	Type          string
	Local         bool
	PHP           string
	AdminUsername string
	AdminPassword string
	AdminEmail    string
	Domain        string
	RootKey       string
	RootCert      string
	SiteCert      string
	SiteKey       string
	Viper         *viper.Viper
}

// loadGlobalConfig gets config information that transcends sites such as app and default settings
func (c *Config) loadGlobalConfig() error {

	dynamicConfig, err := c.loadGlobalViper()
	if err != nil {
		return err
	}

	c.Global.Viper = dynamicConfig
	c.Global.Xdebug = dynamicConfig.GetBool("xdebug")
	c.Global.Local = dynamicConfig.GetBool("local")
	c.Global.AdminEmail = dynamicConfig.GetString("admin.email")
	c.Global.AdminPassword = dynamicConfig.GetString("admin.password")
	c.Global.AdminUsername = dynamicConfig.GetString("admin.username")
	c.Global.PHP = dynamicConfig.GetString("php")
	c.Global.Type = dynamicConfig.GetString("type")

	return err
}

// loadGlobalViper loads the global config using viper and sets defaults
func (c *Config) loadGlobalViper() (*viper.Viper, error) {

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
	dynamicConfig.AddConfigPath(path.Join(c.Directories.App, "config"))

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
	if !isValidString(dynamicConfig.GetString("type"), validTypes) {
		changeConfig = true
		dynamicConfig.Set("type", "site")
	}

	// Reset default php version if there's an invalid version in the config file
	if !isValidString(dynamicConfig.GetString("php"), validPHPVersions) {
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
