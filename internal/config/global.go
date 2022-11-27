package config

import (
	"path"

	"github.com/spf13/viper"
)

// loadGlobalConfig gets config information that transcends sites such as app and default settings
func (s *Settings) loadGlobalConfig() error {

	globalViperConfig, err := s.loadGlobalViper()
	if err != nil {
		return err
	}

	s.global = globalViperConfig
	s.Xdebug = globalViperConfig.GetBool("xdebug")
	s.Local = globalViperConfig.GetBool("local")
	s.AdminEmail = globalViperConfig.GetString("admin.email")
	s.AdminPassword = globalViperConfig.GetString("admin.password")
	s.AdminUsername = globalViperConfig.GetString("admin.username")
	s.PHP = globalViperConfig.GetString("php")
	s.Type = globalViperConfig.GetString("type")

	return err
}

// loadGlobalViper loads the global config using viper and sets defaults
func (s *Settings) loadGlobalViper() (*viper.Viper, error) {

	globalViperConfig := viper.New()

	globalViperConfig.SetDefault("xdebug", xdebug)
	globalViperConfig.SetDefault("type", siteType)
	globalViperConfig.SetDefault("local", local)
	globalViperConfig.SetDefault("php", php)
	globalViperConfig.SetDefault("admin.username", adminUsername)
	globalViperConfig.SetDefault("admin.password", adminPassword)
	globalViperConfig.SetDefault("admin.email", adminEmail)

	globalViperConfig.SetConfigName("kana")
	globalViperConfig.SetConfigType("json")
	globalViperConfig.AddConfigPath(path.Join(s.AppDirectory, "config"))

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
