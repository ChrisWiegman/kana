package config

import (
	"os"
)

type Config struct {
	WorkingDirectory string
	App              AppConfig
	Site             SiteConfig
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

func NewConfig() (Config, error) {

	kanaConfig := Config{}

	cwd, err := os.Getwd()
	if err != nil {
		return kanaConfig, err
	}

	kanaConfig.WorkingDirectory = cwd

	appConfig, err := kanaConfig.getAppConfig()
	if err != nil {
		return kanaConfig, err
	}

	kanaConfig.App = appConfig

	siteConfig, err := kanaConfig.getSiteConfig()
	if err != nil {
		return kanaConfig, err
	}

	kanaConfig.Site = siteConfig

	return kanaConfig, nil

}
