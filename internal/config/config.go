package config

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Directories struct {
	App     string
	Working string
	Site    string
}

type Config struct {
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

var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"
var domain = "sites.kana.li"
var configFolderName = ".config/kana"

func NewConfig() (*Config, error) {

	kanaConfig := new(Config)

	kanaConfig.AppDomain = domain
	kanaConfig.RootKey = rootKey
	kanaConfig.RootCert = rootCert
	kanaConfig.SiteCert = siteCert
	kanaConfig.SiteKey = siteKey

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

	err = kanaConfig.EnsureStaticConfigFiles()
	if err != nil {
		return kanaConfig, err
	}

	err = kanaConfig.loadGlobalConfig()
	if err != nil {
		return kanaConfig, err
	}

	err = kanaConfig.loadLocalConfig()

	return kanaConfig, err
}
