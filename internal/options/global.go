package options

import (
	"errors"
	"path/filepath"

	"github.com/spf13/viper"
)

func loadGlobalOptions(appDirectory string) (*viper.Viper, error) {
	globalSettings := viper.New()

	setViperDefaults(globalSettings, &defaultOptions)

	globalSettings.SetConfigName("kana")
	globalSettings.SetConfigType("json")
	globalSettings.AddConfigPath(filepath.Join(appDirectory, "config"))

	err := globalSettings.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError

		if errors.As(err, &configFileNotFoundError) {
			err = globalSettings.SafeWriteConfig()
			if err != nil {
				return globalSettings, err
			}
		} else {
			return globalSettings, err
		}
	}

	err = validateViperSettings(globalSettings, &defaultOptions)

	return globalSettings, err
}
