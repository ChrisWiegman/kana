package settings

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ChrisWiegman/kana/internal/helpers"
	"github.com/spf13/cobra"
)

func Load(settings *Settings, version string, cmd *cobra.Command, commandsRequiringSite []string) (err error) {
	settings.directories, err = getStaticDirectories()
	if err != nil {
		return err
	}

	settings.constants = appConstants
	settings.constants.Version = version

	err = loadSiteOptions(settings, cmd)
	if err != nil {
		return err
	}

	// Fail now if we have a command that requires a completed site and we haven't started it before
	if settings.site.IsNew && helpers.ArrayContains(commandsRequiringSite, cmd.Use) {
		return fmt.Errorf("the current site you are trying to work with does not exist. Use `kana start` to initialize")
	}

	err = saveLocalLinkConfig(cmd, settings.directories.Site, settings.directories.Working, settings.site.IsNamed)
	if err != nil {
		return err
	}

	settings.directories.Site = filepath.Join(settings.directories.App, "sites", settings.site.Name)

	settings.global, err = loadGlobalOptions(settings.directories.App)
	if err != nil {
		return err
	}

	updateSettingsFromViper(settings.global, &settings.settings)

	settings.local, err = loadLocalOptions(settings.directories.Working, &settings.settings)
	if err != nil {
		return err
	}

	updateSettingsFromViper(settings.local, &settings.settings)

	// Always make sure we set the correct type, even if a config file isn't available.
	if cmd.Use != "start" {
		err = detectType(settings)
		if err != nil {
			return err
		}
	}

	return ensureStaticConfigFiles(settings.directories.App)
}

// DetectType determines the type of site in the working directory.
func detectType(settings *Settings) error {
	var err error
	var isSite bool

	isSite, err = helpers.PathExists(filepath.Join(settings.directories.Working, "wp-includes", "version.php"))
	if err != nil {
		return err
	}

	if isSite {
		return err
	}

	items, _ := os.ReadDir(settings.directories.Working)

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		if item.Name() == "style.css" || filepath.Ext(item.Name()) == ".php" {
			var f *os.File
			var line string

			f, err = os.Open(filepath.Join(settings.directories.Working, item.Name()))
			if err != nil {
				return err
			}

			reader := bufio.NewReader(f)
			line, err = helpers.ReadLine(reader)

			for err == nil {
				exp := regexp.MustCompile(`(Plugin|Theme) Name: .*`)

				for _, match := range exp.FindAllStringSubmatch(line, -1) {
					if match[1] == "Theme" {
						settings.settings.Type = "theme"
					} else {
						settings.settings.Type = "plugin"
					}

					return err
				}
				line, err = helpers.ReadLine(reader)
			}
		}
	}

	// We don't care if it is an empty folder.
	if err == io.EOF {
		err = nil
	}

	settings.settings.Type = "site"

	return err
}

// Get retrieves a setting from the settings object by name and returns it as a string.
// Returns an empty string if the setting is not found.
func (s *Settings) Get(name string) string {
	settingType, settingValue, err := s.getSetting(name)
	if err != nil {
		return ""
	}

	switch settingType {
	case "string":
		return settingValue.String()
	case "bool":
		return strconv.FormatBool(settingValue.Bool())
	case "int64":
		return strconv.FormatInt(settingValue.Int(), 10)
	case "":
		return strings.Join(settingValue.Interface().([]string), ", ")
	default:
		return ""
	}
}

// GetArray retrieves a setting from the settings object by name and returns it as a string array.
func (s *Settings) GetArray(name string) []string {
	return strings.Split(s.Get(name), ", ")
}

// GetBool retrieves a setting from the settings object by name and returns it as a boolean. Returns false if the value is not a boolean.
func (s *Settings) GetBool(name string) bool {
	value, err := strconv.ParseBool(s.Get(name))
	if err != nil {
		return false
	}

	return value
}

// GetGlobalSetting Retrieves a global setting for the "config" command.
func (s *Settings) GetGlobalSetting(name string) (string, error) {
	_, _, err := s.getSetting(name)
	if err != nil {
		return "", fmt.Errorf("invalid setting %s. Please enter a valid key to set", name)
	}

	return s.global.GetString(name), nil
}

// GetInt retrieves a setting from the settings object by name and returns it as an integer. Returns 0 if the value is not an integer.
func (s *Settings) GetInt(name string) int64 {
	value, err := strconv.ParseInt(s.Get(name), 10, 64)
	if err != nil {
		return 0
	}

	return value
}

func (s *Settings) Set(name, value string) error {
	_, _, err := s.getSetting(name)
	if err != nil {
		return fmt.Errorf("invalid setting %s: Please enter a valid key to set", name)
	}

	err = validateSetting(name, value, s.settings.Database)
	if err != nil {
		return err
	}

	name = strings.ToLower(name)

	switch name {
	case "activate",
		"automaticlogin",
		"mailpit",
		"removedefaultplugins",
		"scriptdebug", "ssl",
		"wpdebug",
		"xdebug":
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		s.global.Set(name, boolValue)
	case "updateinterval":
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		s.global.Set(name, intValue)
	default:
		s.global.Set(name, value)
	}

	if name == "database" && value != s.global.GetString("database") {
		switch value {
		case "mariadb":
			s.global.Set("databaseversion", mariadbVersion)
		case "mysql": //nolint:goconst
			s.global.Set("databaseversion", mysqlVersion)
		}
	}

	return s.global.WriteConfig()
}

// getSetting retrieves a setting from the settings object by name and returns the type and value.
// Returns an error if the setting is not found.
func (s *Settings) getSetting(name string) (string, reflect.Value, error) {
	settings := []interface{}{
		&s.settings,
		&s.constants,
		&s.site}

	for group := range settings {
		reflection := reflect.ValueOf(settings[group]).Elem()
		reflectType := reflection.Type()

		for i := 0; i < reflectType.NumField(); i++ {
			field := reflectType.Field(i)
			if name == field.Name || strings.EqualFold(name, field.Name) {
				reflectValue := reflect.ValueOf(settings[group])
				value := reflect.Indirect(reflectValue).FieldByName(field.Name)

				return field.Type.Name(), value, nil
			}
		}
	}

	return "", reflect.Value{}, fmt.Errorf("setting %s not found", name)
}
