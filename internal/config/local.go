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
func (c *Config) ProcessNameFlag(cmd *cobra.Command) (bool, error) {

	isSite := false // Don't assume we're in a site that has been initialized.

	// Don't run this on commands that wouldn't possibly use it.
	if cmd.Use == "config" || cmd.Use == "version" || cmd.Use == "help" {
		return isSite, nil
	}

	// By default the siteLink should be the working directory (assume it's linked)
	siteLink := c.Directories.Working

	// Process the name flag if set
	if cmd.Flags().Lookup("name").Changed {

		// Check that we're not using invalid start flags for the start command
		if cmd.Use == "start" {
			if cmd.Flags().Lookup("plugin").Changed || cmd.Flags().Lookup("theme").Changed || cmd.Flags().Lookup("local").Changed {
				return isSite, fmt.Errorf("invalid flags detected. 'plugin' 'theme' and 'local' flags are not valid with named sites")
			}
		}

		c.Name = sanitizeSiteName(cmd.Flags().Lookup("name").Value.String())
		c.Directories.Site = (path.Join(c.Directories.App, "sites", c.Name))

		c.SiteDomain = fmt.Sprintf("%s.%s", c.Name, c.AppDomain)
		c.SecureURL = fmt.Sprintf("https://%s/", c.SiteDomain)
		c.URL = fmt.Sprintf("http://%s/", c.SiteDomain)

		siteLink = c.Directories.Site
	}

	_, err := os.Stat(path.Join(c.Directories.Site, "link.json"))
	if err == nil || !os.IsNotExist(err) {
		isSite = true
	}

	siteLinkConfig := viper.New()

	siteLinkConfig.SetDefault("link", siteLink)

	siteLinkConfig.SetConfigName("link")
	siteLinkConfig.SetConfigType("json")
	siteLinkConfig.AddConfigPath(c.Directories.Site)

	err = siteLinkConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok && cmd.Use == "start" {
			isSite = true
			err = os.MkdirAll(c.Directories.Site, 0750)
			if err != nil {
				return isSite, err
			}
			err = siteLinkConfig.SafeWriteConfig()
			if err != nil {
				return isSite, err
			}
		}
	}

	c.Directories.Working = siteLinkConfig.GetString("link")

	return isSite, nil
}

// ProcessStartFlags Process the start flags and save them to the settings object
func (c *Config) ProcessStartFlags(cmd *cobra.Command, flags StartFlags) {

	if cmd.Flags().Lookup("local").Changed {
		c.Local = flags.Local
	}

	if cmd.Flags().Lookup("xdebug").Changed {
		c.Xdebug = flags.Local
	}

	if cmd.Flags().Lookup("plugin").Changed && flags.IsPlugin {
		c.Type = "plugih"
	}

	if cmd.Flags().Lookup("theme").Changed && flags.IsTheme {
		c.Type = "theme"
	}
}

// WriteLocalSettings Writes all appropriate local settings to the local config file
func (c *Config) WriteLocalSettings(localSettings LocalSettings) error {

	c.local.Set("local", localSettings.Local)
	c.local.Set("type", localSettings.Type)
	c.local.Set("xdebug", localSettings.Xdebug)
	c.local.Set("plugins", localSettings.Plugins)

	if _, err := os.Stat(path.Join(c.Directories.Working, ".kana.json")); os.IsNotExist(err) {
		return c.local.SafeWriteConfig()
	}

	return c.local.WriteConfig()
}

// loadLocalConfig Loads the config for the current site being called
func (c *Config) loadLocalConfig() error {

	siteName := sanitizeSiteName(filepath.Base(c.Directories.Working))
	// Setup other options generated from config items
	c.SiteDomain = fmt.Sprintf("%s.%s", siteName, c.AppDomain)
	c.SecureURL = fmt.Sprintf("https://%s/", c.SiteDomain)
	c.URL = fmt.Sprintf("http://%s/", c.SiteDomain)

	c.Name = siteName
	c.Directories.Site = path.Join(c.Directories.App, "sites", siteName)

	localViper, err := c.loadlocalViper()
	if err != nil {
		return err
	}

	c.local = localViper
	c.Xdebug = localViper.GetBool("xdebug")
	c.Local = localViper.GetBool("local")
	c.PHP = localViper.GetString("php")
	c.Type = localViper.GetString("type")
	c.Plugins = localViper.GetStringSlice("plugins")

	return nil
}

// loadSiteConfig Get the config items that can be overridden locally with a .kana.json file.
func (c *Config) loadlocalViper() (*viper.Viper, error) {

	localConfig := viper.New()

	localConfig.SetDefault("php", c.PHP)
	localConfig.SetDefault("type", c.Type)
	localConfig.SetDefault("local", c.Local)
	localConfig.SetDefault("xdebug", c.Xdebug)
	localConfig.SetDefault("plugins", []string{})

	localConfig.SetConfigName(".kana")
	localConfig.SetConfigType("json")
	localConfig.AddConfigPath(c.Directories.Working)

	err := localConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return localConfig, err
		}
	}

	return localConfig, nil
}
