package config

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// The following are the default settings for Kana
var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"
var domain = "sites.kana.li"
var configFolderName = ".config/kana"
var php = "8.1"
var siteType = "site"
var xdebug = false
var local = false
var adminUsername = "admin"
var adminPassword = "password"
var adminEmail = "admin@sites.kana.li"

type Settings struct {
	Directories                              Directories
	Local, Xdebug                            bool
	AdminEmail, AdminPassword, AdminUsername string
	AppDomain, SiteDomain                    string
	Name                                     string
	PHP                                      string
	RootCert, RootKey, SiteCert, SiteKey     string
	SecureURL, URL                           string
	Type                                     string
	Plugins                                  []string
	global                                   *viper.Viper
	local                                    *viper.Viper
}

type Directories struct {
	App     string
	Working string
	Site    string
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

func NewConfig() (*Settings, error) {

	kanaSettings := new(Settings)

	kanaSettings.AppDomain = domain
	kanaSettings.RootKey = rootKey
	kanaSettings.RootCert = rootCert
	kanaSettings.SiteCert = siteCert
	kanaSettings.SiteKey = siteKey

	cwd, err := os.Getwd()
	if err != nil {
		return kanaSettings, err
	}

	kanaSettings.Directories.Working = cwd

	home, err := homedir.Dir()
	if err != nil {
		return kanaSettings, err
	}

	kanaSettings.Directories.App = filepath.Join(home, configFolderName)

	err = kanaSettings.EnsureStaticConfigFiles()
	if err != nil {
		return kanaSettings, err
	}

	err = kanaSettings.loadGlobalConfig()
	if err != nil {
		return kanaSettings, err
	}

	err = kanaSettings.loadLocalConfig()

	return kanaSettings, err
}
