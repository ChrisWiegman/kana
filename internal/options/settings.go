package options

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ChrisWiegman/kana/internal/helpers"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func New(version string, cmd *cobra.Command) (*Settings, error) {
	kanaSettings := new(Settings)
	settings := map[string]interface{}{}
	var err error

	for _, setting := range defaults {
		setting.currentValue = setting.defaultValue
		kanaSettings.settings = append(kanaSettings.settings, setting)
	}

	settings["appDirectory"], settings["workingDirectory"], err = getStaticDirectories()
	if err != nil {
		return kanaSettings, err
	}

	settings["name"],
		settings["siteDirectory"],
		settings["isNamed"],
		settings["isNew"],
		err = getSiteInfo(settings["workingDirectory"].(string), settings["appDirectory"].(string), cmd)
	if err != nil {
		return kanaSettings, err
	}

	for key, value := range settings {
		err = kanaSettings.Set(key, value)
		if err != nil {
			return kanaSettings, err
		}
	}

	global, err := getKoanfOptions("global", kanaSettings)
	if err != nil {
		return kanaSettings, err
	}

	kanaSettings.global = global

	local, err := getKoanfOptions("local", kanaSettings)
	if err != nil {
		return kanaSettings, err
	}

	kanaSettings.local = local

	err = ensureStaticConfigFiles(settings["appDirectory"].(string))

	return kanaSettings, err
}

func (s *Settings) Get(name string) string {
	for _, setting := range s.settings {
		if setting.name == name {
			return setting.currentValue
		}
	}

	return ""
}

func (s *Settings) GetBool(name string) bool {
	for _, setting := range s.settings {
		if setting.name == name {
			return setting.currentValue == "true"
		}
	}

	return false
}

func (s *Settings) GetInt(name string) int64 {
	for _, setting := range s.settings {
		if setting.name == name {
			value, err := strconv.ParseInt(setting.currentValue, 10, 64)
			if err != nil {
				return 0
			}

			return value
		}
	}

	return 0
}

func (s *Settings) GetSlice(name string) []string {
	for _, setting := range s.settings {
		if setting.name == name {
			if setting.currentValue == "" {
				return []string{}
			}

			return strings.Split(setting.currentValue, ",")
		}
	}

	return []string{}
}

func (s *Settings) Set(name string, value interface{}) error {
	for i, setting := range s.settings {
		if setting.name == name {
			s.settings[i].currentValue = fmt.Sprint(value)
			return nil
		}
	}

	return fmt.Errorf("invalid setting %s. Please enter a valid key to set", name)
}

func (s *Settings) getAll(settingsType string) map[string]interface{} {
	allSettings := make(map[string]interface{})

	for _, setting := range s.settings {
		if (!setting.hasLocal && settingsType == "local") || (!setting.hasGlobal && settingsType == "global") {
			continue
		}

		switch setting.settingType {
		case "bool":
			boolValue, _ := strconv.ParseBool(setting.currentValue)

			allSettings[setting.name] = boolValue
		case "int":
			intValue, _ := strconv.ParseInt(setting.currentValue, 10, 64)

			allSettings[setting.name] = intValue
		case "slice":
			if setting.currentValue == "" {
				allSettings[setting.name] = []string{}
			} else {
				allSettings[setting.name] = strings.Split(setting.currentValue, ",")
			}
		default:
			allSettings[setting.name] = setting.currentValue
		}
	}

	return allSettings
}

func getSiteInfo(workingDirectory, appDirectory string, cmd *cobra.Command) (name, siteDirectory string, isNamed, isNew bool, err error) {
	name = helpers.SanitizeSiteName(filepath.Base(workingDirectory))
	isStartCommand := cmd.Use == "start"

	// Don't run this on commands that wouldn't possibly use it.
	if cmd.Use == "config" || cmd.Use == "version" || cmd.Use == "help" {
		return name, siteDirectory, isNamed, isNew, err
	}

	// Process the name flag if set
	if cmd.Flags().Lookup("name").Changed {
		isNamed = true

		// Check that we're not using invalid start flags for the start command
		if isStartCommand && cmd.Flags().Lookup("type").Changed {
			typeValue, _ := cmd.Flags().GetString("type")
			if typeValue != "site" {
				return name, siteDirectory, isNamed, isNew,
					fmt.Errorf("the type flag is not valid when using the `name` flag")
			}
		}

		name = helpers.SanitizeSiteName(cmd.Flags().Lookup("name").Value.String())
	}

	// We can set the site directory here now that we have the correct name.
	siteDirectory = filepath.Join(appDirectory, "sites", name)

	_, err = os.Stat(siteDirectory)
	if err != nil && os.IsNotExist(err) {
		if os.IsNotExist(err) {
			isNew = true
		} else {
			return name, siteDirectory, isNamed, isNew, err
		}
	}

	return name, siteDirectory, isNamed, isNew, nil
}

func getStaticDirectories() (app, working string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return app, working, err
	}

	working = cwd

	home, err := homedir.Dir()
	if err != nil {
		return app, working, err
	}

	app = filepath.Join(home, configFolderName)

	err = os.MkdirAll(app, os.FileMode(defaultDirPermissions))

	return app, working, err
}
