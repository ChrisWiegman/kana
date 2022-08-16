package config

import (
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"
var appDomain = "sites.kana.li"
var configFolderName = ".config/kana"

var ValidPHPVersions = []string{
	"site",
	"plugin",
	"theme",
}

var ValidTypes = []string{
	"site",
	"plugin",
	"theme",
}

type AppConfig struct {
	AppDomain            string
	SiteName             string
	AppDirectory         string
	SiteDirectory        string
	RootKey              string
	RootCert             string
	SiteCert             string
	SiteKey              string
	SiteXdebug           bool
	SiteLocal            bool
	SiteType             string
	DefaultPHPVersion    string
	DefaultAdminUsername string
	DefaultAdminPassword string
	DefaultAdminEmail    string
	WorkingDirectory     string
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

	siteName := filepath.Base(cwd)

	viperConfig, err := getViperConfig(appDirectory)
	if err != nil {
		panic(err)
	}

	kanaConfig := AppConfig{
		AppDomain:            appDomain,
		SiteName:             siteName,
		AppDirectory:         appDirectory,
		SiteDirectory:        path.Join(appDirectory, "sites", siteName),
		RootKey:              rootKey,
		RootCert:             rootCert,
		SiteCert:             siteCert,
		SiteKey:              siteKey,
		SiteXdebug:           viperConfig.GetBool("xdebug"),
		SiteLocal:            viperConfig.GetBool("xdebug"),
		SiteType:             viperConfig.GetString("type"),
		DefaultPHPVersion:    viperConfig.GetString("php"),
		DefaultAdminUsername: viperConfig.GetString("adminUser"),
		DefaultAdminPassword: viperConfig.GetString("adminPassword"),
		DefaultAdminEmail:    viperConfig.GetString("adminEmail"),
		WorkingDirectory:     cwd,
	}

	return kanaConfig, nil

}

func getViperConfig(appDirectory string) (*viper.Viper, error) {

	defaultConfig := viper.New()

	defaultConfig.SetDefault("xdebug", false)
	defaultConfig.SetDefault("type", "site")
	defaultConfig.SetDefault("local", false)
	defaultConfig.SetDefault("php", "7.4")
	defaultConfig.SetDefault("adminUser", "admin")
	defaultConfig.SetDefault("adminPassword", "password")
	defaultConfig.SetDefault("adminEmail", "admin@mykanasite.localhost")

	defaultConfig.SetConfigName("kana")
	defaultConfig.SetConfigType("json")
	defaultConfig.AddConfigPath(path.Join(appDirectory, "config"))

	err := defaultConfig.ReadInConfig()
	if err != nil {

		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {

			err = defaultConfig.SafeWriteConfig()
			if err != nil {
				return defaultConfig, err
			}
		} else {
			return defaultConfig, err
		}
	}

	changeConfig := false

	// Reset default "site" type if there's an invalid type in the config file
	if !CheckString(defaultConfig.GetString("type"), ValidTypes) {
		changeConfig = true
		defaultConfig.Set("type", "site")
	}

	// Reset default php version if there's an invalid version in the config file
	if !CheckString(defaultConfig.GetString("php"), ValidPHPVersions) {
		changeConfig = true
		defaultConfig.Set("php", "7.4")
	}

	if changeConfig {
		err = defaultConfig.WriteConfig()
		if err != nil {
			return defaultConfig, err
		}
	}

	return defaultConfig, nil

}

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
