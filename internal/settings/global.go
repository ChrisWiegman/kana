package settings

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func loadGlobalOptions(appDirectory string) (*viper.Viper, error) {
	globalSettings := viper.New()

	setViperDefaults(globalSettings, &defaultOptions)

	configPath := filepath.Join(appDirectory, "config")

	globalSettings.SetConfigName("kana")
	globalSettings.SetConfigType("json")
	globalSettings.AddConfigPath(configPath)

	if err := os.MkdirAll(configPath, os.FileMode(defaultDirPermissions)); err != nil {
		return globalSettings, err
	}

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
