package site

import (
	"github.com/ChrisWiegman/kana/internal/appConfig"

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
