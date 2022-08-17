package appConfig

import (
	"path/filepath"

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

// getAppDirectory Return the path for the global config.
func getAppDirectory() (string, error) {

	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, configFolderName), nil

}
