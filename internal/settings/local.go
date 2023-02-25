package settings

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type StartFlags struct {
	Xdebug     bool
	PhpMyAdmin bool
	Mailpit    bool
	Local      bool
	IsTheme    bool
	IsPlugin   bool
}

type LocalSettings struct {
	Local, PhpMyAdmin, Mailpit, Xdebug, SSL bool
	Type                                    string
	Plugins                                 []string
}

// LoadLocalSettings Loads the config for the current site being called
func (s *Settings) LoadLocalSettings(cmd *cobra.Command) (bool, error) {
	siteName := sanitizeSiteName(filepath.Base(s.WorkingDirectory))
	// Setup other options generated from config items
	s.SiteDomain = fmt.Sprintf("%s.%s", siteName, s.AppDomain)
	s.Protocol = s.getProtocol()
	s.URL = fmt.Sprintf("%s://%s", s.Protocol, s.SiteDomain)

	s.Name = siteName
	s.SiteDirectory = path.Join(s.AppDirectory, "sites", siteName)

	isSite, err := s.ProcessNameFlag(cmd)
	if err != nil {
		return isSite, err
	}

	localViper, err := s.loadlocalViper()
	if err != nil {
		return isSite, err
	}

	s.local = localViper
	s.Xdebug = localViper.GetBool("xdebug")
	s.PhpMyAdmin = localViper.GetBool("phpmyadmin")
	s.Mailpit = localViper.GetBool("mailpit")
	s.Local = localViper.GetBool("local")
	s.PHP = localViper.GetString("php")
	s.Type = localViper.GetString("type")
	s.Plugins = localViper.GetStringSlice("plugins")
	s.SSL = localViper.GetBool("ssl")

	return isSite, nil
}

// ProcessNameFlag Processes the name flag on the site resetting all appropriate local variables
func (s *Settings) ProcessNameFlag(cmd *cobra.Command) (bool, error) {
	isSite := false // Don't assume we're in a site that has been initialized.

	// Don't run this on commands that wouldn't possibly use it.
	if cmd.Use == "config" || cmd.Use == "version" || cmd.Use == "help" {
		return isSite, nil
	}

	// By default the siteLink should be the working directory (assume it's linked)
	siteLink := s.WorkingDirectory

	// Process the name flag if set
	if cmd.Flags().Lookup("name").Changed {
		// Check that we're not using invalid start flags for the start command
		if cmd.Use == "start" {
			if cmd.Flags().Lookup("plugin").Changed || cmd.Flags().Lookup("theme").Changed || cmd.Flags().Lookup("local").Changed {
				return isSite, fmt.Errorf("invalid flags detected. 'plugin' 'theme' and 'local' flags are not valid with named sites")
			}
		}

		s.Name = sanitizeSiteName(cmd.Flags().Lookup("name").Value.String())
		s.SiteDirectory = (path.Join(s.AppDirectory, "sites", s.Name))

		s.SiteDomain = fmt.Sprintf("%s.%s", s.Name, s.AppDomain)
		s.URL = fmt.Sprintf("%s://%s", s.Protocol, s.SiteDomain)

		siteLink = s.SiteDirectory
	}

	_, err := os.Stat(path.Join(s.SiteDirectory, "link.json"))
	if err == nil || !os.IsNotExist(err) {
		isSite = true
	}

	return s.saveLinkConfig(isSite, cmd, siteLink)
}

func (s *Settings) getProtocol() string {
	if s.SSL {
		return "https"
	}

	return "http"
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

// ProcessStartFlags Process the start flags and save them to the settings object
func (s *Settings) ProcessStartFlags(cmd *cobra.Command, flags StartFlags) {
	if cmd.Flags().Lookup("local").Changed {
		s.Local = flags.Local
	}

	if cmd.Flags().Lookup("xdebug").Changed {
		s.Xdebug = flags.Xdebug
	}

	if cmd.Flags().Lookup("phpmyadmin").Changed {
		s.PhpMyAdmin = flags.PhpMyAdmin
	}

	if cmd.Flags().Lookup("mailpit").Changed {
		s.Mailpit = flags.Mailpit
	}

	if cmd.Flags().Lookup("plugin").Changed && flags.IsPlugin {
		s.Type = "plugin"
	}

	if cmd.Flags().Lookup("theme").Changed && flags.IsTheme {
		s.Type = "theme"
	}
}

// WriteLocalSettings Writes all appropriate local settings to the local config file
func (s *Settings) WriteLocalSettings(localSettings LocalSettings) error {
	s.local.Set("local", localSettings.Local)
	s.local.Set("type", localSettings.Type)
	s.local.Set("xdebug", localSettings.Xdebug)
	s.local.Set("phpmyadmin", localSettings.PhpMyAdmin)
	s.local.Set("mailpit", localSettings.Mailpit)
	s.local.Set("plugins", localSettings.Plugins)
	s.local.Set("ssl", localSettings.SSL)

	if _, err := os.Stat(path.Join(s.WorkingDirectory, ".kana.json")); os.IsNotExist(err) {
		return s.local.SafeWriteConfig()
	}

	return s.local.WriteConfig()
}

// loadSiteConfig Get the config items that can be overridden locally with a .kana.json file.
func (s *Settings) loadlocalViper() (*viper.Viper, error) {
	localSettings := viper.New()

	localSettings.SetDefault("php", s.PHP)
	localSettings.SetDefault("type", s.Type)
	localSettings.SetDefault("local", s.Local)
	localSettings.SetDefault("xdebug", s.Xdebug)
	localSettings.SetDefault("phpmyadmin", s.PhpMyAdmin)
	localSettings.SetDefault("mailpit", s.Mailpit)
	localSettings.SetDefault("plugins", []string{})
	localSettings.SetDefault("ssl", s.SSL)

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
