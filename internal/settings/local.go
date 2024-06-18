package settings

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func loadLocalOptions(workingDirectory string, settings *Options) (*viper.Viper, error) {
	localSettings := viper.New()

	setViperDefaults(localSettings, settings)

	localSettings.SetConfigName(".kana")
	localSettings.SetConfigType("json")
	localSettings.AddConfigPath(workingDirectory)

	err := localSettings.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return localSettings, err
		}
	}

	err = validateViperSettings(localSettings, settings)

	return localSettings, err
}

func saveLocalLinkConfig(cmd *cobra.Command, siteDirectory, workingDirectory string, isNamedSite bool) error {
	siteLink := workingDirectory

	if isNamedSite {
		siteLink = siteDirectory
	}

	siteLinkConfig := viper.New()

	siteLinkConfig.SetDefault("link", siteLink)

	siteLinkConfig.SetConfigName("link")
	siteLinkConfig.SetConfigType("json")
	siteLinkConfig.AddConfigPath(siteDirectory)

	err := siteLinkConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok && cmd.Use == "start" { //nolint:goconst
			err = os.MkdirAll(siteDirectory, os.FileMode(defaultDirPermissions))
			if err != nil {
				return err
			}
			err = siteLinkConfig.SafeWriteConfig()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
