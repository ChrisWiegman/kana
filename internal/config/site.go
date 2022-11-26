package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ChrisWiegman/kana-cli/internal/appConfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type StartFlags struct {
	Xdebug   bool
	Local    bool
	IsTheme  bool
	IsPlugin bool
}

type SiteConfig struct {
	SiteName  string
	Domain    string
	URL       string
	SecureURL string
	PHP       string
	Xdebug    bool
	Local     bool
	Type      string
	Plugins   []string
	Viper     *viper.Viper
}

func (c *Config) LoadSiteConfig() error {

	siteName := sanitizeSiteName(filepath.Base(c.Directories.Working))
	// Setup other options generated from config items
	c.Site.Domain = fmt.Sprintf("%s.%s", siteName, c.App.AppDomain)
	c.Site.SecureURL = fmt.Sprintf("https://%s/", c.Site.Domain)
	c.Site.URL = fmt.Sprintf("http://%s/", c.Site.Domain)

	c.Site.SiteName = siteName
	c.Directories.Site = path.Join(c.Directories.App, "sites", siteName)

	siteViper, err := c.loadSiteViper()
	if err != nil {
		return err
	}

	c.Site.Viper = siteViper
	c.Site.Xdebug = siteViper.GetBool("xdebug")
	c.Site.Local = siteViper.GetBool("local")
	c.Site.PHP = siteViper.GetString("php")
	c.Site.Type = siteViper.GetString("type")
	c.Site.Plugins = siteViper.GetStringSlice("plugins")

	return nil
}

// loadSiteConfig Get the config items that can be overridden locally with a .kana.json file.
func (c *Config) loadSiteViper() (*viper.Viper, error) {

	siteConfig := viper.New()

	siteConfig.SetDefault("php", c.App.PHP)
	siteConfig.SetDefault("type", c.App.Type)
	siteConfig.SetDefault("local", c.App.Local)
	siteConfig.SetDefault("xdebug", c.App.Xdebug)
	siteConfig.SetDefault("plugins", []string{})

	siteConfig.SetConfigName(".kana")
	siteConfig.SetConfigType("json")
	siteConfig.AddConfigPath(c.Directories.Working)

	err := siteConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return siteConfig, err
		}
	}

	return siteConfig, nil
}

// ProcessStartFlags Process the start flags and save them to the settings object
func (c *Config) ProcessStartFlags(cmd *cobra.Command, flags StartFlags) {

	if cmd.Flags().Lookup("local").Changed {
		c.Site.Viper.Set("local", flags.Local)
	}

	if cmd.Flags().Lookup("xdebug").Changed {
		c.Site.Viper.Set("xdebug", flags.Xdebug)
	}

	if cmd.Flags().Lookup("plugin").Changed && flags.IsPlugin {
		c.Site.Viper.Set("type", "plugin")
	}

	if cmd.Flags().Lookup("theme").Changed && flags.IsTheme {
		c.Site.Viper.Set("type", "theme")
	}
}

// ProcessNameFlag Processes the name flag on the site resetting all appropriate site variables
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

		c.Site.SiteName = appConfig.SanitizeSiteName(cmd.Flags().Lookup("name").Value.String())
		c.Directories.Site = (path.Join(c.Directories.App, "sites", c.Site.SiteName))

		c.Site.Domain = fmt.Sprintf("%s.%s", c.Site.SiteName, c.App.AppDomain)
		c.Site.SecureURL = fmt.Sprintf("https://%s/", c.Site.Domain)
		c.Site.URL = fmt.Sprintf("http://%s/", c.Site.Domain)

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
