package config

import (
	"os"
	"path"
	"path/filepath"
)

type KanaConfig struct {
	SiteDomain       string
	CurrentDirectory string
	ConfigRoot       string
	SSLCerts         KanaSSLCerts
}

type KanaSSLCerts struct {
	CertDirectory string
	RootKey       string
	RootCert      string
	SiteCert      string
	SiteKey       string
}

var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"

func GetKanaConfig() (KanaConfig, error) {

	configRoot, err := GetConfigRoot()
	if err != nil {
		return KanaConfig{}, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	certDir := path.Join(configRoot, "certs")

	certs := KanaSSLCerts{
		CertDirectory: certDir,
		RootKey:       rootKey,
		RootCert:      rootCert,
		SiteCert:      siteCert,
		SiteKey:       siteKey,
	}

	kanaConfig := KanaConfig{
		SiteDomain:       "sites.cfw.li",
		CurrentDirectory: filepath.Base(cwd),
		ConfigRoot:       configRoot,
		SSLCerts:         certs,
	}

	return kanaConfig, nil

}
