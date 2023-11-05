package settings

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// The following are the default settings for Kana.
var (
	rootKey              = "kana.root.key"
	rootCert             = "kana.root.pem"
	siteCert             = "kana.site.pem"
	siteKey              = "kana.site.key"
	domain               = "sites.kana.li"
	configFolderName     = ".config/kana"
	php                  = "8.1"
	mariadb              = "11"
	siteType             = "site"
	xdebug               = false
	mailpit              = false
	local                = false
	ssl                  = false
	wpdebug              = false
	removeDefaultPlugins = false
	adminUsername        = "admin"
	adminPassword        = "password"
	adminEmail           = "admin@sites.kana.li"
	imageUpdateDays      = 7
	activate             = true
	databaseClient       = "phpmyadmin"
	multisite            = "none"
)

// Default permissions for all new files and folders.
const (
	defaultDirPermissions  = 0750
	defaultFilePermissions = 0644
)

// Settings Individual Settings for use throughout the app lifecycle.
type Settings struct {
	Local, Xdebug, Mailpit, SSL, WPDebug, Activate bool
	RemoveDefaultPlugins                           bool
	ImageUpdateDays                                int
	AdminEmail, AdminPassword, AdminUsername       string
	AppDirectory, SiteDirectory, WorkingDirectory  string
	AppDomain, SiteDomain                          string
	Name                                           string
	PHP                                            string
	MariaDB                                        string
	RootCert, RootKey, SiteCert, SiteKey           string
	URL, Protocol                                  string
	Type                                           string
	DatabaseClient                                 string
	Multisite                                      string
	Plugins                                        []string
	global                                         *viper.Viper
	local                                          *viper.Viper
}

var validPHPVersions = []string{
	"7.4",
	"8.0",
	"8.1",
	"8.2",
}

var validMariaDBVersions = []string{
	"10",
	"11",
}

var validTypes = []string{
	"site",
	"plugin",
	"theme",
}

var validDatabaseClients = []string{
	"phpmyadmin",
	"tableplus",
}

var validMultisiteTypes = []string{
	"none",
	"subdomain",
	"subdirectory",
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

// GetDefaultPermissions returns the default directory permissions and the default file permissions.
func GetDefaultPermissions() (dirPerms, filePerms int) {
	return defaultDirPermissions, defaultFilePermissions
}
