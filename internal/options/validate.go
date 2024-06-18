package options

import (
	"fmt"
	"strings"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/helpers"

	"github.com/go-playground/validator/v10"
)

var validTypes = []string{
	"site",
	"plugin",
	"theme",
}

var validDatabaseClients = []string{
	"phpmyadmin",
	"tableplus",
}

var validDatabases = []string{
	"mariadb",
	"mysql",
	"sqlite",
}

var validMultisiteTypes = []string{
	"none",
	"subdomain",
	"subdirectory",
}

var validEnvironmentTypes = []string{
	"local",
	"development",
	"staging",
	"production",
}

// isValidString Checks a given string against an array of valid values and returns true/false as appropriate.
func isValidString(stringToCheck string, validStrings []string) bool {
	for _, validString := range validStrings {
		if validString == stringToCheck {
			return true
		}
	}

	return false
}

// validateSetting validates the values in saved settings.
func validateSetting(setting, value, database string) error { //nolint:gocyclo
	validate := validator.New()

	setting = strings.ToLower(setting)

	switch setting {
	case "admin.email":
		return validate.Var(value, "email")
	case "admin.password":
		return validate.Var(value, "alphanumunicode")
	case "admin.username":
		return validate.Var(value, "alpha")
	case "database":
		if !helpers.IsValidString(value, validDatabases) {
			return fmt.Errorf("the database, %s, is not a valid database type. You must use either `mariadb`, `mysql` or `sqlite`", setting)
		}
	case "databaseclient":
		if !helpers.IsValidString(value, validDatabaseClients) {
			return fmt.Errorf("the database client, %s, is not a valid client. You must use either `phpmyadmin` or `tableplus`", setting)
		}
	case "environment":
		if !helpers.IsValidString(value, validEnvironmentTypes) {
			return fmt.Errorf("the environment, %s, is not valid. You must use either `local`, `development`, `staging` or `production`", setting)
		}
	case "updateinterval":
		return validate.Var(value, "gte=0")
	case "databaseVersion":
	case "mariadb":
		if docker.ValidateImage(database, value) != nil {
			databaseURL := "https://hub.docker.com/_/mariadb"

			if database == "mysql" {
				databaseURL = "https://hub.docker.com/_/mysql"
			}

			return fmt.Errorf(
				"the database version in your configuration, %s, is invalid. See %s for a list of supported versions",
				value, databaseURL)
		}
	case "multisite":
		if !helpers.IsValidString(value, validMultisiteTypes) {
			return fmt.Errorf("the multisite type, %s, is not a valid type. You must use either `none`, `subdomain` or `subdirectory`", setting)
		}
	case "php":
		if docker.ValidateImage("wordpress", fmt.Sprintf("php%s", value)) != nil {
			return fmt.Errorf(
				"the PHP version in your configuration, %s, is invalid. See https://hub.docker.com/_/wordpress for a list of supported versions",
				value)
		}
	case "plugins":
	case "theme":
		return validate.Var(value, "ascii")
	case "type":
		if !helpers.IsValidString(value, validTypes) {
			return fmt.Errorf("the type you selected, %s, is not a valid type. You must use either `site`, `plugin` or `theme`", setting)
		}
	default:
		err := validate.Var(value, "boolean")
		if err != nil {
			return fmt.Errorf("the setting, %s, must be either true or false", setting)
		}
	}

	return nil
}
