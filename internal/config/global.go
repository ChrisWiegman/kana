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
	Local, Xdebug                            bool
	AdminEmail, AdminPassword, AdminUsername string
	Domain                                   string
	PHP                                      string
	RootCert, RootKey, SiteCert, SiteKey     string
	Type                                     string
	viper                                    *viper.Viper
}

// loadGlobalConfig gets config information that transcends sites such as app and default settings
func (c *Config) loadGlobalConfig() error {

	globalViperConfig, err := c.loadGlobalViper()
	if err != nil {
		return err
	}

	c.Global.viper = globalViperConfig
	c.Global.Xdebug = globalViperConfig.GetBool("xdebug")
	c.Global.Local = globalViperConfig.GetBool("local")
	c.Global.AdminEmail = globalViperConfig.GetString("admin.email")
	c.Global.AdminPassword = globalViperConfig.GetString("admin.password")
	c.Global.AdminUsername = globalViperConfig.GetString("admin.username")
	c.Global.PHP = globalViperConfig.GetString("php")
	c.Global.Type = globalViperConfig.GetString("type")

	return err
}

// loadGlobalViper loads the global config using viper and sets defaults
func (c *Config) loadGlobalViper() (*viper.Viper, error) {

	globalViperConfig := viper.New()

	globalViperConfig.SetDefault("xdebug", false)
	globalViperConfig.SetDefault("type", "site")
	globalViperConfig.SetDefault("local", false)
	globalViperConfig.SetDefault("php", "7.4")
	globalViperConfig.SetDefault("admin.username", "admin")
	globalViperConfig.SetDefault("admin.password", "password")
	globalViperConfig.SetDefault("admin.email", "admin@mykanasite.localhost")

	globalViperConfig.SetConfigName("kana")
	globalViperConfig.SetConfigType("json")
	globalViperConfig.AddConfigPath(path.Join(c.Directories.App, "config"))

	err := globalViperConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			err = globalViperConfig.SafeWriteConfig()
			if err != nil {
				return globalViperConfig, err
			}
		} else {
			return globalViperConfig, err
		}
	}

	changeConfig := false

	// Reset default "site" type if there's an invalid type in the config file
	if !isValidString(globalViperConfig.GetString("type"), validTypes) {
		changeConfig = true
		globalViperConfig.Set("type", "site")
	}

	// Reset default php version if there's an invalid version in the config file
	if !isValidString(globalViperConfig.GetString("php"), validPHPVersions) {
		changeConfig = true
		globalViperConfig.Set("php", "7.4")
	}

	if changeConfig {
		err = globalViperConfig.WriteConfig()
		if err != nil {
			return globalViperConfig, err
		}
	}

	return globalViperConfig, nil
}
