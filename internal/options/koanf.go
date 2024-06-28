package options

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

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

func getKoanfOptions(settingsType string, settings *Settings) (*koanf.Koanf, error) {
	ko := koanf.New(".")

	configFile := getConfigFile(settingsType, settings.Get("workingDirectory"), settings.Get("appDirectory"))
	configFileExists := true

	_, err := os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		configFileExists = false
		if settingsType == "global" { //nolint:goconst
			err = writeKoanfSettings(settingsType, settings)
			if err != nil {
				return ko, err
			}
		}
	}

	if settingsType != "local" || configFileExists {
		err = ko.Load(file.Provider(configFile), kjson.Parser())
		if err != nil {
			return ko, err
		}
	}

	for _, setting := range settings.settings {
		if ko.Exists(setting.name) {
			switch setting.settingType {
			case "bool":
				err = settings.Set(setting.name, ko.Bool(setting.name))
				if err != nil {
					return ko, err
				}
			case "int":
				err = settings.Set(setting.name, ko.Int64(setting.name))
				if err != nil {
					return ko, err
				}
			case "slice":
				stringValue := ko.String(setting.name)
				if stringValue != "" {
					err = settings.Set(setting.name, stringValue)
					if err != nil {
						return ko, err
					}
				} else {
					err = settings.Set(setting.name, strings.Split(setting.currentValue, ","))
					if err != nil {
						return ko, err
					}
				}
			default:
				err = settings.Set(setting.name, ko.String(setting.name))
				if err != nil {
					return ko, err
				}
			}
		}
	}

	return ko, nil
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
