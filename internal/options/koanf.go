package options

import (
	"encoding/json"
	"os"
	"path/filepath"

	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func getConfigFile(settingsType, workingDirectory, appDirectory string) string {
	configFile := filepath.Join(appDirectory, "config", "kana.json")

	if settingsType == "local" { //nolint:goconst
		configFile = filepath.Join(workingDirectory, ".kana.json")
	}

	return configFile
}

func loadKoanfOptions(settingsType string, settings *Settings) error {
	ko := koanf.New(".")

	configFile := getConfigFile(settingsType, settings.Get("workingDirectory"), settings.Get("appDirectory"))
	configFileExists := true

	_, err := os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		configFileExists = false
		if settingsType == "global" { //nolint:goconst
			err = writeKoanfSettings(settingsType, settings)
			if err != nil {
				return err
			}
		}
	}

	if settingsType != "local" || configFileExists {
		err = ko.Load(file.Provider(configFile), kjson.Parser())
		if err != nil {
			return err
		}
	}

	for _, setting := range settings.settings {
		if ko.Exists(setting.name) {
			switch setting.settingType {
			case "bool":
				err = settings.Set(setting.name, ko.Bool(setting.name))
				if err != nil {
					return err
				}
			case "int":
				err = settings.Set(setting.name, ko.Int64(setting.name))
				if err != nil {
					return err
				}
			case "slice":
				err = settings.Set(setting.name, ko.Strings(setting.name))
				if err != nil {
					return err
				}
			default:
				err = settings.Set(setting.name, ko.String(setting.name))
				if err != nil {
					return err
				}
			}
		}
	}

	if settingsType == "local" {
		settings.local = ko
	} else {
		settings.global = ko
	}

	return nil
}

func writeKoanfSettings(settingsType string, settings *Settings) error {
	configFile := getConfigFile(settingsType, settings.Get("workingDirectory"), settings.Get("appDirectory"))
	if settingsType == "global" {
		err := os.MkdirAll(filepath.Dir(configFile), defaultDirPermissions)
		if err != nil {
			return err
		}
	}

	allSettings := settings.getAll(settingsType)

	f, _ := os.Create(configFile)
	defer f.Close()

	jsonBytes, err := json.MarshalIndent(allSettings, "", "\t")
	if err != nil {
		return err
	}

	_, err = f.Write(jsonBytes)

	return err
}
