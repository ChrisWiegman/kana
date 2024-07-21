package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func Load(kanaSettings *Settings, version string, cmd *cobra.Command) error {
	settings := map[string]interface{}{}
	var err error

	for i := range defaults {
		defaults[i].currentValue = defaults[i].defaultValue
		kanaSettings.settings = append(kanaSettings.settings, defaults[i])
	}

	settings["appDirectory"], settings["workingDirectory"], err = getStaticDirectories()
	if err != nil {
		return err
	}

	settings["name"],
		settings["siteDirectory"],
		settings["isNamed"],
		settings["isNew"],
		err = getSiteInfo(settings["workingDirectory"].(string), settings["appDirectory"].(string), cmd)
	if err != nil {
		return err
	}

	for key, value := range settings {
		err = kanaSettings.Set(key, value)
		if err != nil {
			return err
		}
	}

	err = saveLocalLinkConfig(cmd, settings["siteDirectory"].(string), settings["workingDirectory"].(string), settings["isNamed"].(bool))
	if err != nil {
		return err
	}

	err = loadKoanfOptions("global", kanaSettings)
	if err != nil {
		return err
	}

	err = loadKoanfOptions("local", kanaSettings)
	if err != nil {
		return err
	}

	err = ensureStaticConfigFiles(settings["appDirectory"].(string))
	if err != nil {
		return err
	}

	return processStartFlags(cmd, kanaSettings)
}

func (s *Settings) Get(name string) string {
	for i := range s.settings {
		if strings.EqualFold(s.settings[i].name, name) {
			return s.settings[i].currentValue
		}
	}

	switch name {
	case "rootCert":
		return rootCert
	case "rootKey":
		return rootKey
	case "siteCert":
		return siteCert
	case "siteKey":
		return siteKey
	}

	return ""
}

func (s *Settings) GetAll(settingsType string) map[string]interface{} {
	allSettings := make(map[string]interface{})
	koSettings := s.global

	if settingsType == "local" {
		koSettings = s.local
	}

	for i := range s.settings {
		if (!s.settings[i].hasLocal && settingsType == "local") || (!s.settings[i].hasGlobal && settingsType == "global") {
			continue
		}

		switch s.settings[i].settingType {
		case "bool":
			boolValue, _ := strconv.ParseBool(s.settings[i].currentValue)
			if koSettings != nil && koSettings.Exists(s.settings[i].name) {
				boolValue = koSettings.Bool(s.settings[i].name)
			}

			allSettings[s.settings[i].name] = boolValue
		case "int":
			intValue, _ := strconv.ParseInt(s.settings[i].currentValue, 10, 64)
			if koSettings != nil && koSettings.Exists(s.settings[i].name) {
				intValue = koSettings.Int64(s.settings[i].name)
			}

			allSettings[s.settings[i].name] = intValue
		case "slice":
			sliceVal := strings.Split(s.settings[i].currentValue, ",")
			if koSettings != nil && koSettings.Exists(s.settings[i].name) {
				sliceVal = koSettings.Strings(s.settings[i].name)
			}

			allSettings[s.settings[i].name] = sliceVal
		default:
			stringValue := s.settings[i].currentValue
			if koSettings != nil && koSettings.Exists(s.settings[i].name) {
				stringValue = koSettings.String(s.settings[i].name)
			}

			allSettings[s.settings[i].name] = stringValue
		}
	}

	return allSettings
}

func (s *Settings) GetBool(name string) bool {
	for i := range s.settings {
		if strings.EqualFold(s.settings[i].name, name) {
			return s.settings[i].currentValue == "true"
		}
	}

	return false
}

func (s *Settings) GetInt(name string) int64 {
	for i := range s.settings {
		if strings.EqualFold(s.settings[i].name, name) {
			value, err := strconv.ParseInt(s.settings[i].currentValue, 10, 64)
			if err != nil {
				return 0
			}

			return value
		}
	}

	return 0
}

func (s *Settings) GetSlice(name string) []string {
	for i := range s.settings {
		if strings.EqualFold(s.settings[i].name, name) {
			if s.settings[i].currentValue == "" {
				return []string{}
			}

			return strings.Split(s.settings[i].currentValue, ",")
		}
	}

	return []string{}
}

func (s *Settings) Set(name string, value interface{}, setVars ...bool) error {
	for i := range s.settings {
		if !strings.EqualFold(s.settings[i].name, name) {
			continue
		}

		err := s.validate(name, value)
		if err != nil {
			return err
		}

		if s.settings[i].settingType == "slice" && reflect.TypeOf(value).String() == "[]string" {
			s.settings[i].currentValue = strings.Join(value.([]string), ",")
		} else {
			s.settings[i].currentValue = fmt.Sprint(value)
		}

		if len(setVars) > 0 && setVars[0] {
			if s.settings[i].settingType == "slice" {
				value = strings.Split(s.settings[i].currentValue, ",")
			}

			err := s.global.Set(s.settings[i].name, value)
			if err != nil {
				return err
			}

			err = writeKoanfSettings("global", s)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("invalid setting %s. Please enter a valid key to set", name)
}

func (s *Settings) WriteLocalSettings(localSettings map[string]interface{}) error {
	configFile := filepath.Join(s.Get("workingDirectory"), ".kana.json")

	allSettings := s.GetAll("local")

	for setting, value := range localSettings {
		allSettings[setting] = value
	}

	f, _ := os.Create(configFile)
	defer f.Close()

	jsonBytes, err := json.MarshalIndent(allSettings, "", "\t")
	if err != nil {
		return err
	}

	_, err = f.Write(jsonBytes)

	return err
}

func (s *Settings) validate(name string, value interface{}) error {
	for i := range s.settings {
		if !strings.EqualFold(s.settings[i].name, name) {
			continue
		}

		stringVal := fmt.Sprint(value)

		switch s.settings[i].settingType {
		case "bool":
			_, err := strconv.ParseBool(stringVal)
			if err != nil {
				return fmt.Errorf("the value for %s must be a boolean", name)
			}
		case "int":
			_, err := strconv.ParseInt(stringVal, 10, 64)
			if err != nil {
				return fmt.Errorf("the value for %s must be an integer", name)
			}
		}

		if len(s.settings[i].validValues) > 0 {
			if !helpers.IsValidString(stringVal, s.settings[i].validValues) {
				return fmt.Errorf("the %s value, %s, is not valid", name, stringVal)
			}
		}

		validate := validator.New()

		switch name {
		case "adminEmail":
			return validate.Var(stringVal, "email")
		case "updateInterval":
			return validate.Var(stringVal, "gte=0")
		case "databaseVersion":
			if docker.ValidateImage(s.Get("database"), stringVal) != nil {
				databaseURL := "https://hub.docker.com/_/mariadb"

				if s.Get("database") == "mysql" {
					databaseURL = "https://hub.docker.com/_/mysql"
				}

				return fmt.Errorf(
					"the database version in your configuration, %s, is invalid. See %s for a list of supported versions",
					stringVal, databaseURL)
			}
		case "php":
			if docker.ValidateImage("wordpress", fmt.Sprintf("php%s", stringVal)) != nil {
				return fmt.Errorf(
					"the PHP version in your configuration, %s, is invalid. See https://hub.docker.com/_/wordpress for a list of supported versions",
					stringVal)
			}
		}
	}

	return nil
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

func saveLocalLinkConfig(cmd *cobra.Command, siteDirectory, workingDirectory string, isNamedSite bool) error {
	siteLink := map[string]string{
		"link": workingDirectory}

	if isNamedSite {
		siteLink["link"] = siteDirectory
	}

	linkConfigFile := filepath.Join(siteDirectory, "link.json")

	_, err := os.Stat(linkConfigFile)

	if err != nil && os.IsNotExist(err) && cmd.Use == "start" {
		err := os.MkdirAll(filepath.Dir(linkConfigFile), defaultDirPermissions)
		if err != nil {
			return err
		}

		f, _ := os.Create(linkConfigFile)
		defer f.Close()

		jsonBytes, err := json.MarshalIndent(siteLink, "", "\t")
		if err != nil {
			return err
		}

		_, err = f.Write(jsonBytes)
		if err != nil {
			return err
		}
	}

	return nil
}
