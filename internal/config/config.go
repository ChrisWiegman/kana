package config

import (
	"fmt"
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

type AppConfig struct {
	AppDomain     string
	SiteName      string
	AppDirectory  string
	SiteDirectory string
	RootKey       string
	RootCert      string
	SiteCert      string
	SiteKey       string
	SiteXdebug    bool
	SiteLocal     bool
	SiteType      string
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
		AppDomain:     appDomain,
		SiteName:      siteName,
		AppDirectory:  appDirectory,
		SiteDirectory: path.Join(appDirectory, "sites", siteName),
		RootKey:       rootKey,
		RootCert:      rootCert,
		SiteCert:      siteCert,
		SiteKey:       siteKey,
		SiteXdebug:    viperConfig.GetBool("xdebug"),
		SiteLocal:     viperConfig.GetBool("xdebug"),
		SiteType:      viperConfig.GetString("type"),
	}

	return kanaConfig, nil

}

func getViperConfig(appDirectory string) (*viper.Viper, error) {

	defaultConfig := viper.New()

	defaultConfig.SetDefault("xdebug", false)
	defaultConfig.SetDefault("type", "site")
	defaultConfig.SetDefault("local", false)

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

	if !checkValidType(defaultConfig.GetString("type")) {
		return viper.New(), fmt.Errorf("the value of site in your app config is not a valid config. possible values are [site, plugin, theme]")
	}

	return defaultConfig, nil

}

func checkValidType(typeToCheck string) bool {

	validTypes := []string{
		"site",
		"plugin",
		"theme",
	}

	for _, validType := range validTypes {
		if validType == typeToCheck {
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
