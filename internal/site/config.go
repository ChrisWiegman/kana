package site

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ChrisWiegman/kana-cli/internal/appConfig"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SiteFlags struct {
	Xdebug   bool
	Local    bool
	IsTheme  bool
	IsPlugin bool
}

// getSiteConfig Get the config items that can be overridden locally with a .kana.json file.
func getSiteConfig(staticConfig appConfig.StaticConfig, dynamicConfig *viper.Viper) (*viper.Viper, error) {

	siteConfig := viper.New()

	siteConfig.SetDefault("php", dynamicConfig.GetString("php"))
	siteConfig.SetDefault("type", dynamicConfig.GetString("type"))
	siteConfig.SetDefault("local", dynamicConfig.GetBool("local"))
	siteConfig.SetDefault("xdebug", dynamicConfig.GetBool("xdebug"))
	siteConfig.SetDefault("plugins", []string{})

	siteConfig.SetConfigName(".kana")
	siteConfig.SetConfigType("json")
	siteConfig.AddConfigPath(staticConfig.WorkingDirectory)

	err := siteConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return siteConfig, err
		}
	}

	return siteConfig, nil
}

func (s *Site) ExportSiteConfig() error {

	config := s.GetRunningConfig()
	plugins, err := s.GetInstalledWordPressPlugins()
	if err != nil {
		return err
	}

	s.SiteConfig.Set("local", config.Local)
	s.SiteConfig.Set("type", config.Type)
	s.SiteConfig.Set("xdebug", config.Xdebug)
	s.SiteConfig.Set("plugins", plugins)

	if _, err = os.Stat(path.Join(s.StaticConfig.WorkingDirectory, ".kana.json")); os.IsNotExist(err) {
		return s.SiteConfig.SafeWriteConfig()
	}

	return s.SiteConfig.WriteConfig()
}

// IsLocalSite Determines if a site is a "local" site (started with the "local" flag) so that other commands can work as needed.
func (s *Site) IsLocalSite() bool {

	// If the site is already running, try to make this easier
	if s.IsSiteRunning() {
		runningConfig := s.GetRunningConfig()
		if runningConfig.Local {
			return true
		}
	}

	// First check the app site folders. If we've created the site (has a DB) without an "app" folder we can assume it was local last time.
	hasNonLocalAppFolder := true
	hasDatabaseFolder := true

	if _, err := os.Stat(path.Join(s.StaticConfig.SiteDirectory, "app")); os.IsNotExist(err) {
		hasNonLocalAppFolder = false
	}

	if _, err := os.Stat(path.Join(s.StaticConfig.SiteDirectory, "database")); os.IsNotExist(err) {
		hasDatabaseFolder = false
	}

	if hasDatabaseFolder && !hasNonLocalAppFolder {
		return true
	}

	// Return the flag for all other conditions
	return s.SiteConfig.GetBool("local")
}

// ProcessSiteFlags Process the start flags and save them to the settings object
func (s *Site) ProcessSiteFlags(cmd *cobra.Command, flags SiteFlags) {

	if cmd.Flags().Lookup("local").Changed {
		s.SiteConfig.Set("local", flags.Local)
	}

	if cmd.Flags().Lookup("xdebug").Changed {
		s.SiteConfig.Set("xdebug", flags.Xdebug)
	}

	if cmd.Flags().Lookup("plugin").Changed && flags.IsPlugin {
		s.SiteConfig.Set("type", "plugin")
	}

	if cmd.Flags().Lookup("theme").Changed && flags.IsTheme {
		s.SiteConfig.Set("type", "theme")
	}
}

// GetRunningConfig gets various options that were used to start the site
func (s *Site) GetRunningConfig() CurrentConfig {

	currentConfig := CurrentConfig{
		Type:   "site",
		Local:  false,
		Xdebug: false,
	}

	output, _ := s.runCli("pecl list | grep xdebug", false)
	if strings.Contains(output.StdOut, "xdebug") {
		currentConfig.Xdebug = true
	}

	mounts := s.dockerClient.ContainerGetMounts(fmt.Sprintf("kana_%s_wordpress", s.StaticConfig.SiteName))

	if len(mounts) == 1 {
		currentConfig.Type = "site"
	}

	for _, mount := range mounts {

		if mount.Source == path.Join(s.StaticConfig.WorkingDirectory, "wordpress") {
			currentConfig.Local = true
		}

		if strings.Contains(mount.Destination, "/var/www/html/wp-content/plugins/") {
			currentConfig.Type = "plugin"
		}

		if strings.Contains(mount.Destination, "/var/www/html/wp-content/themes/") {
			currentConfig.Type = "theme"
		}
	}

	return currentConfig
}
