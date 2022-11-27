package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type StartFlags struct {
	Xdebug   bool
	Local    bool
	IsTheme  bool
	IsPlugin bool
}

type LocalSettings struct {
	Local, Xdebug bool
	Type          string
	Plugins       []string
}

// ProcessNameFlag Processes the name flag on the site resetting all appropriate local variables
func (s *Settings) ProcessNameFlag(cmd *cobra.Command) (bool, error) {

	isSite := false // Don't assume we're in a site that has been initialized.

	// Don't run this on commands that wouldn't possibly use it.
	if cmd.Use == "config" || cmd.Use == "version" || cmd.Use == "help" {
		return isSite, nil
	}

	// By default the siteLink should be the working directory (assume it's linked)
	siteLink := s.Directories.Working

	// Process the name flag if set
	if cmd.Flags().Lookup("name").Changed {

		// Check that we're not using invalid start flags for the start command
		if cmd.Use == "start" {
			if cmd.Flags().Lookup("plugin").Changed || cmd.Flags().Lookup("theme").Changed || cmd.Flags().Lookup("local").Changed {
				return isSite, fmt.Errorf("invalid flags detected. 'plugin' 'theme' and 'local' flags are not valid with named sites")
			}
		}

		s.Name = sanitizeSiteName(cmd.Flags().Lookup("name").Value.String())
		s.Directories.Site = (path.Join(s.Directories.App, "sites", s.Name))

		s.SiteDomain = fmt.Sprintf("%s.%s", s.Name, s.AppDomain)
		s.SecureURL = fmt.Sprintf("https://%s/", s.SiteDomain)
		s.URL = fmt.Sprintf("http://%s/", s.SiteDomain)

		siteLink = s.Directories.Site
	}

	_, err := os.Stat(path.Join(s.Directories.Site, "link.json"))
	if err == nil || !os.IsNotExist(err) {
		isSite = true
	}

	siteLinkConfig := viper.New()

	siteLinkConfig.SetDefault("link", siteLink)

	siteLinkConfig.SetConfigName("link")
	siteLinkConfig.SetConfigType("json")
	siteLinkConfig.AddConfigPath(s.Directories.Site)

	err = siteLinkConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok && cmd.Use == "start" {
			isSite = true
			err = os.MkdirAll(s.Directories.Site, 0750)
			if err != nil {
				return isSite, err
			}
			err = siteLinkConfig.SafeWriteConfig()
			if err != nil {
				return isSite, err
			}
		}
	}

	s.Directories.Working = siteLinkConfig.GetString("link")

	return isSite, nil
}

// ProcessStartFlags Process the start flags and save them to the settings object
func (s *Settings) ProcessStartFlags(cmd *cobra.Command, flags StartFlags) {

	if cmd.Flags().Lookup("local").Changed {
		s.Local = flags.Local
	}

	if cmd.Flags().Lookup("xdebug").Changed {
		s.Xdebug = flags.Local
	}

	if cmd.Flags().Lookup("plugin").Changed && flags.IsPlugin {
		s.Type = "plugih"
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
	s.local.Set("plugins", localSettings.Plugins)

	if _, err := os.Stat(path.Join(s.Directories.Working, ".kana.json")); os.IsNotExist(err) {
		return s.local.SafeWriteConfig()
	}

	return s.local.WriteConfig()
}

// loadLocalConfig Loads the config for the current site being called
func (s *Settings) loadLocalConfig() error {

	siteName := sanitizeSiteName(filepath.Base(s.Directories.Working))
	// Setup other options generated from config items
	s.SiteDomain = fmt.Sprintf("%s.%s", siteName, s.AppDomain)
	s.SecureURL = fmt.Sprintf("https://%s/", s.SiteDomain)
	s.URL = fmt.Sprintf("http://%s/", s.SiteDomain)

	s.Name = siteName
	s.Directories.Site = path.Join(s.Directories.App, "sites", siteName)

	localViper, err := s.loadlocalViper()
	if err != nil {
		return err
	}

	s.local = localViper
	s.Xdebug = localViper.GetBool("xdebug")
	s.Local = localViper.GetBool("local")
	s.PHP = localViper.GetString("php")
	s.Type = localViper.GetString("type")
	s.Plugins = localViper.GetStringSlice("plugins")

	return nil
}

// loadSiteConfig Get the config items that can be overridden locally with a .kana.json file.
func (s *Settings) loadlocalViper() (*viper.Viper, error) {

	localConfig := viper.New()

	localConfig.SetDefault("php", s.PHP)
	localConfig.SetDefault("type", s.Type)
	localConfig.SetDefault("local", s.Local)
	localConfig.SetDefault("xdebug", s.Xdebug)
	localConfig.SetDefault("plugins", []string{})

	localConfig.SetConfigName(".kana")
	localConfig.SetConfigType("json")
	localConfig.AddConfigPath(s.Directories.Working)

	err := localConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return localConfig, err
		}
	}

	return localConfig, nil
}
