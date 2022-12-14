package settings

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// The following are the default settings for Kana
var (
	rootKey          = "kana.root.key"
	rootCert         = "kana.root.pem"
	siteCert         = "kana.site.pem"
	siteKey          = "kana.site.key"
	domain           = "sites.kana.li"
	configFolderName = ".config/kana"
	php              = "8.1"
	siteType         = "site"
	xdebug           = false
	phpmyadmin       = false
	local            = false
	adminUsername    = "admin"
	adminPassword    = "password"
	adminEmail       = "admin@sites.kana.li"
)

// Individual Settings for use throughout the app lifecycle
type Settings struct {
	Local, PhpMyAdmin, Xdebug                     bool
	AdminEmail, AdminPassword, AdminUsername      string
	AppDirectory, SiteDirectory, WorkingDirectory string
	AppDomain, SiteDomain                         string
	Name                                          string
	PHP                                           string
	RootCert, RootKey, SiteCert, SiteKey          string
	SecureURL, URL                                string
	Type                                          string
	Plugins                                       []string
	global                                        *viper.Viper
	local                                         *viper.Viper
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

func NewSettings() (*Settings, error) {

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

	kanaSettings.WorkingDirectory = cwd

	home, err := homedir.Dir()
	if err != nil {
		return kanaSettings, err
	}

	kanaSettings.AppDirectory = filepath.Join(home, configFolderName)

	err = kanaSettings.EnsureStaticConfigFiles()

	return kanaSettings, err
}
