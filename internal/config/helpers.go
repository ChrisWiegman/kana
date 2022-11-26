package config

import (
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func CheckString(stringToCheck string, validStrings []string) bool {

	for _, validString := range validStrings {
		if validString == stringToCheck {
			return true
		}
	}

	return false
}

// sanitizeSiteName Returns the site name, properly sanitized for use.
func sanitizeSiteName(rawSiteName string) string {

	siteName := strings.TrimSpace(rawSiteName)
	siteName = strings.ToLower(siteName)
	siteName = strings.ReplaceAll(siteName, " ", "-")
	return strings.ToValidUTF8(siteName, "")
}

// getAppDirectory Return the path for the global config.
func getAppDirectory() (string, error) {

	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, configFolderName), nil
}
