package options

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/docker"

	"github.com/spf13/viper"
)

func getSettingsVromViper(viperSettings *viper.Viper) Options {
	options := Options{}

	options.Activate = viperSettings.GetBool("activate")
	options.AdminEmail = viperSettings.GetString("admin.email")
	options.AdminPassword = viperSettings.GetString("admin.password")
	options.AdminUsername = viperSettings.GetString("admin.username")
	options.AutomaticLogin = viperSettings.GetBool("automaticlogin")
	options.Database = viperSettings.GetString("database")
	options.DatabaseClient = viperSettings.GetString("databaseclient")
	options.DatabaseVersion = viperSettings.GetString("databaseversion")
	options.Environment = viperSettings.GetString("environment")
	options.UpdateInterval = viperSettings.GetInt64("imageupdatedays")
	options.UpdateInterval = viperSettings.GetInt64("updateinterval")
	options.Mailpit = viperSettings.GetBool("mailpit")
	options.Multisite = viperSettings.GetString("multisite")
	options.PHP = viperSettings.GetString("php")
	options.Plugins = viperSettings.GetStringSlice("plugins")
	options.RemoveDefaultPlugins = viperSettings.GetBool("removedefaultplugins")
	options.ScriptDebug = viperSettings.GetBool("scriptdebug")
	options.SSL = viperSettings.GetBool("ssl")
	options.Theme = viperSettings.GetString("theme")
	options.Type = viperSettings.GetString("type")
	options.WPDebug = viperSettings.GetBool("wpdebug")
	options.Xdebug = viperSettings.GetBool("xdebug")

	return options
}

func updateSettingsFromViper(viperSettings *viper.Viper, options *Options) {
	options.Activate = viperSettings.GetBool("activate")
	options.AdminEmail = viperSettings.GetString("admin.email")
	options.AdminPassword = viperSettings.GetString("admin.password")
	options.AdminUsername = viperSettings.GetString("admin.username")
	options.AutomaticLogin = viperSettings.GetBool("automaticlogin")
	options.Database = viperSettings.GetString("database")
	options.DatabaseClient = viperSettings.GetString("databaseclient")
	options.DatabaseVersion = viperSettings.GetString("databaseversion")
	options.Environment = viperSettings.GetString("environment")
	options.UpdateInterval = viperSettings.GetInt64("updateinterval")
	options.Mailpit = viperSettings.GetBool("mailpit")
	options.Multisite = viperSettings.GetString("multisite")
	options.PHP = viperSettings.GetString("php")
	options.Plugins = viperSettings.GetStringSlice("plugins")
	options.RemoveDefaultPlugins = viperSettings.GetBool("removedefaultplugins")
	options.ScriptDebug = viperSettings.GetBool("scriptdebug")
	options.SSL = viperSettings.GetBool("ssl")
	options.Theme = viperSettings.GetString("theme")
	options.Type = viperSettings.GetString("type")
	options.WPDebug = viperSettings.GetBool("wpdebug")
	options.Xdebug = viperSettings.GetBool("xdebug")
}

func setViperDefaults(viperSettings *viper.Viper, options *Options) {
	viperSettings.SetDefault("activate", options.Activate)
	viperSettings.SetDefault("admin.email", options.AdminEmail)
	viperSettings.SetDefault("admin.password", options.AdminPassword)
	viperSettings.SetDefault("admin.username", options.AdminUsername)
	viperSettings.SetDefault("automaticLogin", options.AutomaticLogin)
	viperSettings.SetDefault("database", options.Database)
	viperSettings.SetDefault("databaseclient", options.DatabaseClient)
	viperSettings.SetDefault("databaseversion", options.DatabaseVersion)
	viperSettings.SetDefault("environment", options.Environment)
	viperSettings.SetDefault("updateinterval", options.UpdateInterval)
	viperSettings.SetDefault("mailpit", options.Mailpit)
	viperSettings.SetDefault("multisite", options.Multisite)
	viperSettings.SetDefault("php", options.PHP)
	viperSettings.SetDefault("plugins", options.Plugins)
	viperSettings.SetDefault("removedefaultplugins", options.RemoveDefaultPlugins)
	viperSettings.SetDefault("scriptdebug", options.ScriptDebug)
	viperSettings.SetDefault("ssl", options.SSL)
	viperSettings.SetDefault("theme", options.Theme)
	viperSettings.SetDefault("type", options.Type)
	viperSettings.SetDefault("wpdebug", options.WPDebug)
	viperSettings.SetDefault("xdebug", options.Xdebug)
}

func validateViperSettings(viperSettings *viper.Viper, options *Options) error {
	changeConfig := false

	// Reset default php version if there's an invalid version in the config file
	if docker.ValidateImage("wordpress", fmt.Sprintf("php%s", viperSettings.GetString("php"))) != nil {
		changeConfig = true
		viperSettings.Set("php", options.PHP)
	}

	// Reset default "site" type if there's an invalid type in the config file
	if !isValidString(viperSettings.GetString("type"), validTypes) {
		changeConfig = true
		viperSettings.Set("type", options.Type)
	}

	// Reset default database version if there's an invalid version in the config file
	if docker.ValidateImage(viperSettings.GetString("database"), viperSettings.GetString("databaseVersion")) != nil {
		changeConfig = true
		defaultDatabaseVersion := mariadbVersion

		if viperSettings.GetString("database") == "mysql" {
			defaultDatabaseVersion = mysqlVersion
		}

		viperSettings.Set("databaseVersion", defaultDatabaseVersion)
	}

	// Reset default database type if there's an invalid type in the config file
	if !isValidString(viperSettings.GetString("database"), validDatabases) {
		changeConfig = true
		viperSettings.Set("database", options.Database)
	}

	// Reset default database client if there's an invalid client in the config file
	if !isValidString(viperSettings.GetString("databaseClient"), validDatabaseClients) {
		changeConfig = true
		viperSettings.Set("databaseClient", options.DatabaseClient)
	}

	// Reset default multisite type if there's an invalid type in the config file
	if !isValidString(viperSettings.GetString("multisite"), validMultisiteTypes) {
		changeConfig = true
		viperSettings.Set("multisite", options.Multisite)
	}

	// Reset default environment type if there's an invalid type in the config file
	if !isValidString(viperSettings.GetString("environment"), validEnvironmentTypes) {
		changeConfig = true
		viperSettings.Set("environment", options.Environment)
	}

	if changeConfig {
		err := viperSettings.WriteConfig()
		if err != nil {
			return err
		}
	}

	return nil
}
