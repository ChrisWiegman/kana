package config

import (
	"os"
	"path/filepath"
)

type KanaConfig struct {
	AppDomain     string
	SiteDirectory string
	AppDirectory  string
	RootKey       string
	RootCert      string
	SiteCert      string
	SiteKey       string
}

var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"
var appDomain = "sites.cfw.li"

func GetKanaConfig() (KanaConfig, error) {

	appDirectory, err := GetAppDirectory()
	if err != nil {
		return KanaConfig{}, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	kanaConfig := KanaConfig{
		AppDomain:     appDomain,
		SiteDirectory: filepath.Base(cwd),
		AppDirectory:  appDirectory,
		RootKey:       rootKey,
		RootCert:      rootCert,
		SiteCert:      siteCert,
		SiteKey:       siteKey,
	}

	return kanaConfig, nil

}
