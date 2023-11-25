package settings

import (
	"errors"
	"fmt"
	"path"

	"github.com/ChrisWiegman/kana-cli/internal/docker"
	"github.com/spf13/viper"
)

// LoadGlobalSettings gets config information that transcends sites such as app and default settings.
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
	s.RemoveDefaultPlugins = globalViperConfig.GetBool("removeDefaultPlugins")
	s.AdminEmail = globalViperConfig.GetString("admin.email")
	s.AdminPassword = globalViperConfig.GetString("admin.password")
	s.AdminUsername = globalViperConfig.GetString("admin.username")
	s.PHP = globalViperConfig.GetString("php")
	s.MariaDB = globalViperConfig.GetString("mariadb")
	s.Type = globalViperConfig.GetString("type")
	s.SSL = globalViperConfig.GetBool("ssl")
	s.ImageUpdateDays = globalViperConfig.GetInt("imageUpdateDays")
	s.Activate = globalViperConfig.GetBool("activate")
	s.DatabaseClient = globalViperConfig.GetString("databaseClient")
	s.Multisite = globalViperConfig.GetString("multisite")

	return err
}

// loadGlobalViper loads the global config using viper and sets defaults.
func (s *Settings) loadGlobalViper() (*viper.Viper, error) {
	globalSettings := viper.New()

	globalSettings.SetDefault("xdebug", xdebug)
	globalSettings.SetDefault("imageUpdateDays", imageUpdateDays)
	globalSettings.SetDefault("databaseClient", databaseClient)
	globalSettings.SetDefault("wpdebug", wpdebug)
	globalSettings.SetDefault("mailpit", mailpit)
	globalSettings.SetDefault("type", siteType)
	globalSettings.SetDefault("local", local)
	globalSettings.SetDefault("ssl", ssl)
	globalSettings.SetDefault("activate", activate)
	globalSettings.SetDefault("php", php)
	globalSettings.SetDefault("mariadb", mariadb)
	globalSettings.SetDefault("admin.username", adminUsername)
	globalSettings.SetDefault("admin.password", adminPassword)
	globalSettings.SetDefault("admin.email", adminEmail)
	globalSettings.SetDefault("multisite", multisite)
	globalSettings.SetDefault("removeDefaultPlugins", removeDefaultPlugins)

	globalSettings.SetConfigName("kana")
	globalSettings.SetConfigType("json")
	globalSettings.AddConfigPath(path.Join(s.AppDirectory, "config"))

	err := globalSettings.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError

		if errors.As(err, &configFileNotFoundError) {
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
	if docker.ValidateImage("wordpress", fmt.Sprintf("php%s", globalSettings.GetString("php"))) != nil {
		changeConfig = true
		globalSettings.Set("php", php)
	}

	// Reset default mariadb version if there's an invalid version in the config file
	if docker.ValidateImage("mariadb", globalSettings.GetString("mariadb")) != nil {
		changeConfig = true
		globalSettings.Set("mariadb", mariadb)
	}

	// Reset default database client if there's an invalid client in the config file
	if !isValidString(globalSettings.GetString("databaseClient"), validDatabaseClients) {
		changeConfig = true
		globalSettings.Set("databaseClient", databaseClient)
	}

	// Reset default multisite type if there's an invalid type in the config file
	if !isValidString(globalSettings.GetString("multisite"), validMultisiteTypes) {
		changeConfig = true
		globalSettings.Set("multisite", multisite)
	}

	if changeConfig {
		err = globalSettings.WriteConfig()
		if err != nil {
			return globalSettings, err
		}
	}

	return globalSettings, nil
}
