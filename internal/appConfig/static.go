package appConfig

import (
	"os"
	"path"
	"path/filepath"
)

var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"
var appDomain = "sites.kana.li"
var configFolderName = ".config/kana"

type StaticConfig struct {
	AppDomain        string
	SiteName         string
	AppDirectory     string
	SiteDirectory    string
	RootKey          string
	RootCert         string
	SiteCert         string
	SiteKey          string
	WorkingDirectory string
}

func GetStaticConfig() (StaticConfig, error) {

	appDirectory, err := getAppDirectory()
	if err != nil {
		return StaticConfig{}, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return StaticConfig{}, err
	}

	siteName := filepath.Base(cwd)

	staticConfig := StaticConfig{
		AppDomain:        appDomain,
		SiteName:         siteName,
		AppDirectory:     appDirectory,
		SiteDirectory:    path.Join(appDirectory, "sites", siteName),
		RootKey:          rootKey,
		RootCert:         rootCert,
		SiteCert:         siteCert,
		SiteKey:          siteKey,
		WorkingDirectory: cwd,
	}

	return staticConfig, nil

}
