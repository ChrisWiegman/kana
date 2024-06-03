package settings

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/helpers"
)

// LoadGlobalSettings gets config information that transcends sites such as app and default settings.
func (s *Settings) LoadGlobalSettings() error {
	globalViperConfig, err := s.loadGlobalViper()
	if err != nil {
		return err
	}

	s.global = globalViperConfig

	s.Activate = globalViperConfig.GetBool("activate")
	s.AdminEmail = globalViperConfig.GetString("admin.email")
	s.AdminPassword = globalViperConfig.GetString("admin.password")
	s.AdminUsername = globalViperConfig.GetString("admin.username")
	s.AutomaticLogin = globalViperConfig.GetBool("automaticLogin")
	s.Database = globalViperConfig.GetString("database")
	s.DatabaseClient = globalViperConfig.GetString("databaseClient")
	s.DatabaseVersion = globalViperConfig.GetString("databaseVersion")
	s.Environment = globalViperConfig.GetString("environment")
	s.ImageUpdateDays = globalViperConfig.GetInt("imageUpdateDays")
	s.Mailpit = globalViperConfig.GetBool("mailpit")
	s.Multisite = globalViperConfig.GetString("multisite")
	s.PHP = globalViperConfig.GetString("php")
	s.Plugins = globalViperConfig.GetStringSlice("plugins")
	s.RemoveDefaultPlugins = globalViperConfig.GetBool("removeDefaultPlugins")
	s.ScriptDebug = globalViperConfig.GetBool("scriptdebug")
	s.SSL = globalViperConfig.GetBool("ssl")
	s.Theme = globalViperConfig.GetString(("theme"))
	s.Type = globalViperConfig.GetString("type")
	s.WPDebug = globalViperConfig.GetBool("wpdebug")
	s.Xdebug = globalViperConfig.GetBool("xdebug")

	return err
}

// loadGlobalViper loads the global config using viper and sets defaults.
func (s *Settings) loadGlobalViper() (*viper.Viper, error) { //nolint:funlen
	globalSettings := viper.New()

	globalSettings.SetDefault("activate", activate)
	globalSettings.SetDefault("admin.email", adminEmail)
	globalSettings.SetDefault("admin.password", adminPassword)
	globalSettings.SetDefault("admin.username", adminUsername)
	globalSettings.SetDefault("automaticLogin", automaticLogin)
	globalSettings.SetDefault("database", database)
	globalSettings.SetDefault("databaseClient", databaseClient)
	globalSettings.SetDefault("databaseVersion", databaseVersion)
	globalSettings.SetDefault("environment", environment)
	globalSettings.SetDefault("imageUpdateDays", imageUpdateDays)
	globalSettings.SetDefault("mailpit", mailpit)
	globalSettings.SetDefault("multisite", multisite)
	globalSettings.SetDefault("php", php)
	globalSettings.SetDefault("plugins", plugins)
	globalSettings.SetDefault("removeDefaultPlugins", removeDefaultPlugins)
	globalSettings.SetDefault("scriptdebug", scriptDebug)
	globalSettings.SetDefault("ssl", ssl)
	globalSettings.SetDefault("theme", theme)
	globalSettings.SetDefault("type", siteType)
	globalSettings.SetDefault("wpdebug", wpdebug)
	globalSettings.SetDefault("xdebug", xdebug)

	globalSettings.SetConfigName("kana")
	globalSettings.SetConfigType("json")
	globalSettings.AddConfigPath(filepath.Join(s.AppDirectory, "config"))

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

	// Reset default php version if there's an invalid version in the config file
	if docker.ValidateImage("wordpress", fmt.Sprintf("php%s", globalSettings.GetString("php"))) != nil {
		changeConfig = true
		globalSettings.Set("php", php)
	}

	// Reset default "site" type if there's an invalid type in the config file
	if !helpers.IsValidString(globalSettings.GetString("type"), validTypes) {
		changeConfig = true
		globalSettings.Set("type", "site")
	}

	// Reset default mariadb version if there's an invalid version in the config file
	if docker.ValidateImage(globalSettings.GetString("database"), globalSettings.GetString("databaseVersion")) != nil {
		changeConfig = true
		defaultDatabaseVersion := "11"

		if globalSettings.GetString("database") == "mysql" {
			defaultDatabaseVersion = "8"
		}

		globalSettings.Set("databaseVersion", defaultDatabaseVersion)
	}

	// Reset default database type if there's an invalid type in the config file
	if !helpers.IsValidString(globalSettings.GetString("database"), validDatabases) {
		changeConfig = true
		globalSettings.Set("database", database)
	}

	// Reset default database client if there's an invalid client in the config file
	if !helpers.IsValidString(globalSettings.GetString("databaseClient"), validDatabaseClients) {
		changeConfig = true
		globalSettings.Set("databaseClient", databaseClient)
	}

	// Reset default multisite type if there's an invalid type in the config file
	if !helpers.IsValidString(globalSettings.GetString("multisite"), validMultisiteTypes) {
		changeConfig = true
		globalSettings.Set("multisite", multisite)
	}

	// Reset default environment type if there's an invalid type in the config file
	if !helpers.IsValidString(globalSettings.GetString("environment"), validEnvironmentTypes) {
		changeConfig = true
		globalSettings.Set("environment", environment)
	}

	if changeConfig {
		err = globalSettings.WriteConfig()
		if err != nil {
			return globalSettings, err
		}
	}

	return globalSettings, nil
}
