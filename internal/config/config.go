package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type KanaConfig struct {
	SiteDomain       string
	CurrentDirectory string
	ConfigRoot       string
	SSLCerts         KanaSSLCerts
	HTTPClient       http.Client
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
		RootKey:       path.Join(certDir, rootKey),
		RootCert:      path.Join(certDir, rootCert),
		SiteCert:      path.Join(certDir, siteCert),
		SiteKey:       path.Join(certDir, siteKey),
	}

	kanaConfig := KanaConfig{
		SiteDomain:       "sites.cfw.li",
		CurrentDirectory: filepath.Base(cwd),
		ConfigRoot:       configRoot,
		SSLCerts:         certs,
		HTTPClient:       getSecureHTTPClient(certs),
	}

	return kanaConfig, nil

}

func getSecureHTTPClient(certs KanaSSLCerts) http.Client {

	caCert, err := ioutil.ReadFile(certs.RootCert)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	return client
}
