package config

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type Directories struct {
	App     string
	Working string
	Site    string
}

type Config struct {
	Directories Directories
	App         AppConfig
	Site        SiteConfig
}

var validPHPVersions = []string{
	"7.4",
	"8.0",
	"8.1",
}

var validTypes = []string{
	"site",
	"plugin",
	"theme",
}

func NewConfig() (*Config, error) {

	kanaConfig := new(Config)

	kanaConfig.App = AppConfig{
		Domain:   domain,
		RootKey:  rootKey,
		RootCert: rootCert,
		SiteCert: siteCert,
		SiteKey:  siteKey,
	}

	cwd, err := os.Getwd()
	if err != nil {
		return kanaConfig, err
	}

	kanaConfig.Directories.Working = cwd

	home, err := homedir.Dir()
	if err != nil {
		return kanaConfig, err
	}

	kanaConfig.Directories.App = filepath.Join(home, configFolderName)

	return kanaConfig, nil

}
