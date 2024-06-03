package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChrisWiegman/kana/internal/helpers"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoadLocalSettings Loads the config for the current site being called.
func (s *Settings) LoadLocalSettings(cmd *cobra.Command) (bool, error) {
	siteName := helpers.SanitizeSiteName(filepath.Base(s.WorkingDirectory))
	// Setup other options generated from config items
	s.SiteDomain = fmt.Sprintf("%s.%s", siteName, s.AppDomain)

	s.Name = siteName
	s.SiteDirectory = filepath.Join(s.AppDirectory, "sites", siteName)

	isSite, err := s.ProcessNameFlag(cmd)
	if err != nil {
		return isSite, err
	}

	localViper, err := s.loadLocalViper()
	if err != nil {
		return isSite, err
	}

	s.local = localViper

	s.Activate = localViper.GetBool("activate")
	s.AdminEmail = localViper.GetString("admin.email")
	s.AdminPassword = localViper.GetString("admin.password")
	s.AdminUsername = localViper.GetString("admin.username")
	s.AutomaticLogin = localViper.GetBool("automaticLogin")
	s.Database = localViper.GetString("database")
	s.DatabaseClient = localViper.GetString("databaseClient")
	s.Environment = localViper.GetString("environment")
	s.Mailpit = localViper.GetBool("mailpit")
	s.MariaDB = localViper.GetString("mariadb")
	s.Multisite = localViper.GetString("multisite")
	s.PHP = localViper.GetString("php")
	s.Plugins = localViper.GetStringSlice("plugins")
	s.RemoveDefaultPlugins = localViper.GetBool("removeDefaultPlugins")
	s.ScriptDebug = localViper.GetBool("scriptdebug")
	s.SSL = localViper.GetBool("ssl")
	s.Theme = localViper.GetString("theme")
	s.Type = localViper.GetString("type")
	s.WPDebug = localViper.GetBool("wpdebug")
	s.Xdebug = localViper.GetBool("xdebug")

	s.Protocol = s.getProtocol()
	s.URL = fmt.Sprintf("%s://%s", s.Protocol, s.SiteDomain)

	return isSite, nil
}

// HasLocalOptions Returns true if local options have been saved to a file or false.
func (s *Settings) HasLocalOptions() bool {
	if _, err := os.Stat(filepath.Join(s.WorkingDirectory, ".kana.json")); os.IsNotExist(err) {
		return false
	}

	return true
}

// ProcessNameFlag Processes the name flag on the site resetting all appropriate local variables.
func (s *Settings) ProcessNameFlag(cmd *cobra.Command) (bool, error) { //nolint:gocyclo
	isStartCommand := cmd.Use == "start"

	// Don't run this on commands that wouldn't possibly use it.
	if cmd.Use == "config" || cmd.Use == "version" || cmd.Use == "help" {
		return false, nil
	}

	// Validate the flags that could make this an invalid site
	if isStartCommand && cmd.Flags().Lookup("multisite").Changed {
		multisiteValue, err := cmd.Flags().GetString("multisite")
		if !helpers.IsValidString(multisiteValue, validMultisiteTypes) || err != nil {
			return false,
				fmt.Errorf("the multisite type, %s, is not a valid type. You must use either `none`, `subdomain` or `subdirectory`", multisiteValue)
		}
	}

	if isStartCommand && cmd.Flags().Lookup("type").Changed {
		typeValue, err := cmd.Flags().GetString("type")
		if !helpers.IsValidString(typeValue, validTypes) || err != nil {
			return false,
				fmt.Errorf("the type, %s, is not valid. Only a `site` is valid with named sites", typeValue)
		}
	}

	if isStartCommand && cmd.Flags().Lookup("environment").Changed {
		environmentValue, err := cmd.Flags().GetString("environment")
		if !helpers.IsValidString(environmentValue, validEnvironmentTypes) || err != nil {
			return false,
				fmt.Errorf("the environment, %s, is not valid. You must use either `local`, `development`, `staging` or `production`", environmentValue)
		}
	}

	// By default, the siteLink should be the working directory (assume it's linked)
	siteLink := s.WorkingDirectory

	// Process the name flag if set
	if cmd.Flags().Lookup("name").Changed {
		s.IsNamedSite = true

		// Check that we're not using invalid start flags for the start command
		if isStartCommand && cmd.Flags().Lookup("type").Changed {
			typeValue, _ := cmd.Flags().GetString("type")
			if typeValue != "site" {
				return false,
					fmt.Errorf("the type, %s, is not valid. Only a `site` is valid with named sites", typeValue)
			}
		}

		s.Name = helpers.SanitizeSiteName(cmd.Flags().Lookup("name").Value.String())
		s.SiteDirectory = (filepath.Join(s.AppDirectory, "sites", s.Name))

		s.SiteDomain = fmt.Sprintf("%s.%s", s.Name, s.AppDomain)
		s.URL = fmt.Sprintf("%s://%s", s.Protocol, s.SiteDomain)

		siteLink = s.SiteDirectory
	}

	isSite := false // Don't assume we're in a site that has been initialized.

	_, err := os.Stat(filepath.Join(s.SiteDirectory, "link.json"))
	if err == nil || !os.IsNotExist(err) {
		isSite = true
	}

	s.IsNewSite = !isSite // Negate the site exists here to set if this is a new site.

	return s.saveLinkConfig(isSite, cmd, siteLink)
}

// ProcessStartFlags Process the start flags and save them to the settings object.
func (s *Settings) ProcessStartFlags(cmd *cobra.Command, flags *StartFlags) {
	if cmd.Flags().Lookup("xdebug").Changed {
		s.Xdebug = flags.Xdebug
	}

	if cmd.Flags().Lookup("wpdebug").Changed {
		s.WPDebug = flags.WPDebug
	}

	if cmd.Flags().Lookup("scriptdebug").Changed {
		s.ScriptDebug = flags.ScriptDebug
	}

	if cmd.Flags().Lookup("ssl").Changed {
		s.SSL = flags.SSL
		s.Protocol = s.getProtocol()
		s.URL = fmt.Sprintf("%s://%s", s.Protocol, s.SiteDomain)
	}

	if cmd.Flags().Lookup("mailpit").Changed {
		s.Mailpit = flags.Mailpit
	}

	if cmd.Flags().Lookup("type").Changed {
		s.Type = flags.Type
	}

	if cmd.Flags().Lookup("theme").Changed {
		s.Theme = flags.Theme
	}

	if cmd.Flags().Lookup("multisite").Changed {
		s.Multisite = flags.Multisite
	}

	if cmd.Flags().Lookup("activate").Changed {
		s.Activate = flags.Activate
	}

	if cmd.Flags().Lookup("environment").Changed {
		s.Environment = flags.Environment
	}

	if cmd.Flags().Lookup("plugins").Changed {
		s.Plugins = strings.Split(flags.Plugins, ",")
	}

	if cmd.Flags().Lookup("remove-default-plugins").Changed {
		s.RemoveDefaultPlugins = flags.RemoveDefaultPlugins
	}

	if cmd.Flags().Lookup("database").Changed {
		s.Database = flags.Database
	}
}

// WriteLocalSettings Writes all appropriate local settings to the local config file.
func (s *Settings) WriteLocalSettings(localSettings *LocalSettings) error {
	s.local.Set("activate", localSettings.Activate)
	s.local.Set("admin.email", s.AdminEmail)
	s.local.Set("admin.username", s.AdminUsername)
	s.local.Set("admin.password", s.AdminPassword)
	s.local.Set("automaticLogin", s.AutomaticLogin)
	s.local.Set("database", localSettings.Database)
	s.local.Set("databaseClient", localSettings.DatabaseClient)
	s.local.Set("environment", localSettings.Environment)
	s.local.Set("mailpit", localSettings.Mailpit)
	s.local.Set("multisite", localSettings.Multisite)
	s.local.Set("plugins", localSettings.Plugins)
	s.local.Set("removeDefaultPlugins", localSettings.RemoveDefaultPlugins)
	s.local.Set("scriptdebug", localSettings.ScriptDebug)
	s.local.Set("ssl", localSettings.SSL)
	s.local.Set("theme", localSettings.Theme)
	s.local.Set("type", localSettings.Type)
	s.local.Set("wpdebug", localSettings.WPDebug)
	s.local.Set("xdebug", localSettings.Xdebug)

	if _, err := os.Stat(filepath.Join(s.WorkingDirectory, ".kana.json")); os.IsNotExist(err) {
		return s.local.SafeWriteConfig()
	}

	return s.local.WriteConfig()
}

func (s *Settings) getProtocol() string {
	if s.SSL {
		return "https"
	}

	return "http"
}

// loadLocalViper Get the config items that can be overridden locally with a .kana.json file.
func (s *Settings) loadLocalViper() (*viper.Viper, error) {
	localSettings := viper.New()

	localSettings.SetDefault("activate", s.Activate)
	localSettings.SetDefault("admin.email", s.AdminEmail)
	localSettings.SetDefault("admin.username", s.AdminUsername)
	localSettings.SetDefault("admin.password", s.AdminPassword)
	localSettings.SetDefault("automaticLogin", s.AutomaticLogin)
	localSettings.SetDefault("database", s.Database)
	localSettings.SetDefault("databaseClient", s.DatabaseClient)
	localSettings.SetDefault("environment", s.Environment)
	localSettings.SetDefault("mailpit", s.Mailpit)
	localSettings.SetDefault("mariadb", s.MariaDB)
	localSettings.SetDefault("multisite", s.Multisite)
	localSettings.SetDefault("php", s.PHP)
	localSettings.SetDefault("plugins", s.Plugins)
	localSettings.SetDefault("removeDefaultPlugins", s.RemoveDefaultPlugins)
	localSettings.SetDefault("scriptdebug", s.ScriptDebug)
	localSettings.SetDefault("ssl", s.SSL)
	localSettings.SetDefault("theme", s.Theme)
	localSettings.SetDefault("type", s.Type)
	localSettings.SetDefault("wpdebug", s.WPDebug)
	localSettings.SetDefault("xdebug", s.Xdebug)

	localSettings.SetConfigName(".kana")
	localSettings.SetConfigType("json")
	localSettings.AddConfigPath(s.WorkingDirectory)

	err := localSettings.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return localSettings, err
		}
	}

	return localSettings, nil
}

func (s *Settings) saveLinkConfig(isSite bool, cmd *cobra.Command, siteLink string) (bool, error) {
	siteLinkConfig := viper.New()

	siteLinkConfig.SetDefault("link", siteLink)

	siteLinkConfig.SetConfigName("link")
	siteLinkConfig.SetConfigType("json")
	siteLinkConfig.AddConfigPath(s.SiteDirectory)

	err := siteLinkConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok && cmd.Use == "start" {
			isSite = true
			err = os.MkdirAll(s.SiteDirectory, os.FileMode(defaultDirPermissions))
			if err != nil {
				return isSite, err
			}
			err = siteLinkConfig.SafeWriteConfig()
			if err != nil {
				return isSite, err
			}
		}
	}

	s.WorkingDirectory = siteLinkConfig.GetString("link")

	return isSite, nil
}
