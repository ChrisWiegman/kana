package config

import (
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

type SiteConfig struct {
	SiteName      string
	SiteDirectory string
	PHP           string
	Xdebug        bool
	Local         bool
	Type          string
	Plugins       []string
	Viper         *viper.Viper
}

func (c *Config) LoadSiteConfig() error {

	siteName := sanitizeSiteName(filepath.Base(c.Directories.Working))

	c.Site.SiteName = siteName
	c.Site.SiteDirectory = path.Join(c.Directories.App, "sites", siteName)

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
