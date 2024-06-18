package settings

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func NewSettings(version string, cmd *cobra.Command) (settings *Settings, err error) {
	settings = new(Settings)

	settings.directories, err = getStaticDirectories()
	if err != nil {
		return settings, err
	}

	settings.constants = appConstants
	settings.constants.Version = version

	err = loadSiteOptions(settings, cmd)
	if err != nil {
		return settings, err
	}

	err = saveLocalLinkConfig(cmd, settings.directories.Site, settings.directories.Working, settings.site.IsNamed)
	if err != nil {
		return settings, err
	}

	settings.directories.Site = filepath.Join(settings.directories.App, "sites", settings.site.Name)

	settings.global, err = loadGlobalOptions(settings.directories.App)
	if err != nil {
		return settings, err
	}

	updateSettingsFromViper(settings.global, &settings.settings)

	settings.local, err = loadLocalOptions(settings.directories.Working, &settings.settings)
	if err != nil {
		return settings, err
	}

	updateSettingsFromViper(settings.local, &settings.settings)

	err = ensureStaticConfigFiles(settings.directories.App)

	return settings, err
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
