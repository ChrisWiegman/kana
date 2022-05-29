package config

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var ConfigFolderName = ".kana"

// GetConfigRoot Return the path for the global config.
func GetConfigRoot() (string, error) {

	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ConfigFolderName), nil

}
