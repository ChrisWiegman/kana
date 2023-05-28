package settings

import (
	"path"

	"github.com/spf13/viper"
)

// LoadGlobalSettings gets config information that transcends sites such as app and default settings
func (s *Settings) LoadGlobalSettings() error {
	globalViperConfig, err := s.loadGlobalViper()
	if err != nil {
		return err
	}

	s.global = globalViperConfig
	s.Xdebug = globalViperConfig.GetBool("xdebug")
	s.WPDebug = globalViperConfig.GetBool("wpdebug")
	s.Mailpit = globalViperConfig.GetBool("mailpit")
	s.Local = globalViperConfig.GetBool("local")
	s.AdminEmail = globalViperConfig.GetString("admin.email")
	s.AdminPassword = globalViperConfig.GetString("admin.password")
	s.AdminUsername = globalViperConfig.GetString("admin.username")
	s.PHP = globalViperConfig.GetString("php")
	s.Type = globalViperConfig.GetString("type")
	s.SSL = globalViperConfig.GetBool("ssl")
	s.ImageUpdateDays = globalViperConfig.GetInt("imageUpdateDays")

	return err
}

// loadGlobalViper loads the global config using viper and sets defaults
func (s *Settings) loadGlobalViper() (*viper.Viper, error) {
	globalSettings := viper.New()

	globalSettings.SetDefault("xdebug", xdebug)
	globalSettings.SetDefault("imageUpdateDays", imageUpdateDays)
	globalSettings.SetDefault("wpdebug", wpdebug)
	globalSettings.SetDefault("mailpit", mailpit)
	globalSettings.SetDefault("type", siteType)
	globalSettings.SetDefault("local", local)
	globalSettings.SetDefault("ssl", ssl)
	globalSettings.SetDefault("activate", activate)
	globalSettings.SetDefault("php", php)
	globalSettings.SetDefault("admin.username", adminUsername)
	globalSettings.SetDefault("admin.password", adminPassword)
	globalSettings.SetDefault("admin.email", adminEmail)

	globalSettings.SetConfigName("kana")
	globalSettings.SetConfigType("json")
	globalSettings.AddConfigPath(path.Join(s.AppDirectory, "config"))

	err := globalSettings.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			err = globalSettings.SafeWriteConfig()
			if err != nil {
				return globalSettings, err
			}
		} else {
			return globalSettings, err
		}
	}

	changeConfig := false

	// Reset default "site" type if there's an invalid type in the config file
	if !isValidString(globalSettings.GetString("type"), validTypes) {
		changeConfig = true
		globalSettings.Set("type", "site")
	}

	// Reset default php version if there's an invalid version in the config file
	if !isValidString(globalSettings.GetString("php"), validPHPVersions) {
		changeConfig = true
		globalSettings.Set("php", "7.4")
	}

	if changeConfig {
		err = globalSettings.WriteConfig()
		if err != nil {
			return globalSettings, err
		}
	}

	return globalSettings, nil
}
