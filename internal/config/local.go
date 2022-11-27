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

type LocalConfig struct {
	Name      string
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

		c.Local.Name = sanitizeSiteName(cmd.Flags().Lookup("name").Value.String())
		c.Directories.Site = (path.Join(c.Directories.App, "sites", c.Local.Name))

		c.Local.Domain = fmt.Sprintf("%s.%s", c.Local.Name, c.Global.Domain)
		c.Local.SecureURL = fmt.Sprintf("https://%s/", c.Local.Domain)
		c.Local.URL = fmt.Sprintf("http://%s/", c.Local.Domain)

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
		c.Local.Local = flags.Local
	}

	if cmd.Flags().Lookup("xdebug").Changed {
		c.Local.Xdebug = flags.Local
	}

	if cmd.Flags().Lookup("plugin").Changed && flags.IsPlugin {
		c.Local.Type = "plugih"
	}

	if cmd.Flags().Lookup("theme").Changed && flags.IsTheme {
		c.Local.Type = "theme"
	}
}

// loadLocalConfig Loads the config for the current site being called
func (c *Config) loadLocalConfig() error {

	siteName := sanitizeSiteName(filepath.Base(c.Directories.Working))
	// Setup other options generated from config items
	c.Local.Domain = fmt.Sprintf("%s.%s", siteName, c.Global.Domain)
	c.Local.SecureURL = fmt.Sprintf("https://%s/", c.Local.Domain)
	c.Local.URL = fmt.Sprintf("http://%s/", c.Local.Domain)

	c.Local.Name = siteName
	c.Directories.Site = path.Join(c.Directories.App, "sites", siteName)

	siteViper, err := c.loadSiteViper()
	if err != nil {
		return err
	}

	c.Local.Viper = siteViper
	c.Local.Xdebug = siteViper.GetBool("xdebug")
	c.Local.Local = siteViper.GetBool("local")
	c.Local.PHP = siteViper.GetString("php")
	c.Local.Type = siteViper.GetString("type")
	c.Local.Plugins = siteViper.GetStringSlice("plugins")

	return nil
}

// loadSiteConfig Get the config items that can be overridden locally with a .kana.json file.
func (c *Config) loadSiteViper() (*viper.Viper, error) {

	siteConfig := viper.New()

	siteConfig.SetDefault("php", c.Global.PHP)
	siteConfig.SetDefault("type", c.Global.Type)
	siteConfig.SetDefault("local", c.Global.Local)
	siteConfig.SetDefault("xdebug", c.Global.Xdebug)
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
