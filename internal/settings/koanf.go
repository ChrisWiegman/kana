package settings

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

	for i := range settings.settings {
		if ko.Exists(settings.settings[i].name) {
			switch settings.settings[i].settingType {
			case "bool":
				err = settings.Set(settings.settings[i].name, ko.Bool(settings.settings[i].name))
				if err != nil {
					return err
				}
			case "int": //nolint:goconst
				err = settings.Set(settings.settings[i].name, ko.Int64(settings.settings[i].name))
				if err != nil {
					return err
				}
			case "slice":
				err = settings.Set(settings.settings[i].name, ko.Strings(settings.settings[i].name))
				if err != nil {
					return err
				}
			default:
				err = settings.Set(settings.settings[i].name, ko.String(settings.settings[i].name))
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

	allSettings := settings.GetAll(settingsType)

	f, _ := os.Create(configFile)
	defer f.Close()

	jsonBytes, err := json.MarshalIndent(allSettings, "", "\t")
	if err != nil {
		return err
	}

	_, err = f.Write(jsonBytes)

	return err
}
