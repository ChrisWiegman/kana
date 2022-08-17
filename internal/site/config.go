package site

import (
	"github.com/ChrisWiegman/kana/internal/appConfig"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SiteConfig struct {
	PHPVersion string
	Xdebug     bool
	Local      bool
	Type       string
}

type SiteFlags struct {
	Xdebug   bool
	Local    bool
	IsTheme  bool
	IsPlugin bool
}

func getSiteConfig(staticConfig appConfig.StaticConfig, dynamicConfig appConfig.DynamicConfig) (SiteConfig, error) {

	viperConfig := viper.New()

	viperConfig.SetDefault("php", dynamicConfig.PHPVersion)
	viperConfig.SetDefault("type", dynamicConfig.SiteType)
	viperConfig.SetDefault("local", dynamicConfig.SiteLocal)
	viperConfig.SetDefault("xdebug", dynamicConfig.SiteXdebug)

	viperConfig.SetConfigName(".kana")
	viperConfig.SetConfigType("json")
	viperConfig.AddConfigPath(staticConfig.WorkingDirectory)

	err := viperConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return SiteConfig{}, err
		}
	}

	siteConfig := SiteConfig{
		PHPVersion: viperConfig.GetString("php"),
		Type:       viperConfig.GetString("type"),
		Local:      viperConfig.GetBool("local"),
		Xdebug:     viperConfig.GetBool("xdebug"),
	}

	return siteConfig, nil

}

func (s *Site) ProcessSiteFlags(cmd *cobra.Command, flags SiteFlags) {

	if cmd.Flags().Lookup("local").Changed {
		s.SiteConfig.Local = flags.Local
	}

	if cmd.Flags().Lookup("xdebug").Changed {
		s.SiteConfig.Xdebug = flags.Xdebug
	}

	if cmd.Flags().Lookup("plugin").Changed && flags.IsPlugin {
		s.SiteConfig.Type = "plugin"
	}

	if cmd.Flags().Lookup("theme").Changed && flags.IsTheme {
		s.SiteConfig.Type = "theme"
	}
}
