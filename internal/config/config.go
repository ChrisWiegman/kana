package config

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"
var appDomain = "sites.kana.li"
var configFolderName = ".config/kana"

type AppConfig struct {
	AppDomain    string
	SiteName     string
	AppDirectory string
	RootKey      string
	RootCert     string
	SiteCert     string
	SiteKey      string
}

func GetAppConfig() (AppConfig, error) {

	appDirectory, err := getAppDirectory()
	if err != nil {
		return AppConfig{}, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	kanaConfig := AppConfig{
		AppDomain:    appDomain,
		SiteName:     filepath.Base(cwd),
		AppDirectory: appDirectory,
		RootKey:      rootKey,
		RootCert:     rootCert,
		SiteCert:     siteCert,
		SiteKey:      siteKey,
	}

	return kanaConfig, nil

}

// getAppDirectory Return the path for the global config.
func getAppDirectory() (string, error) {

	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, configFolderName), nil

}
