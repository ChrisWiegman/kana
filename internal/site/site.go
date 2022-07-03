package site

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/pkg/browser"
)

type KanaSite struct {
	rootCert   string
	siteDomain string
	secureURL  string
	url        string
}

func NewSite(config config.KanaConfig) *KanaSite {

	site := new(KanaSite)

	site.rootCert = config.SSLCerts.RootCert
	site.siteDomain = fmt.Sprintf("%s.%s", config.CurrentDirectory, config.SiteDomain)
	site.secureURL = fmt.Sprintf("https://%s/", site.siteDomain)
	site.url = fmt.Sprintf("http://%s/", site.siteDomain)

	return site
}

func (s *KanaSite) GetURL(insecure bool) string {

	if insecure {
		return s.url
	}

	return s.secureURL

}

func (s *KanaSite) OpenSite() error {

	caCert, err := ioutil.ReadFile(s.rootCert)
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	resp, err := client.Get(s.secureURL)
	if err != nil {
		return err
	}

	tries := 0

	for resp.StatusCode != 200 {

		resp, err = client.Get(s.secureURL)
		if err != nil {
			return err
		}

		if resp.StatusCode == 200 {
			break
		}

		if tries == 30 {
			return fmt.Errorf("timeout reached. unable to open site")
		}

		tries++
		time.Sleep(1 * time.Second)

	}

	openURL(s.secureURL)

	return nil

}

func openURL(url string) error {

	if runtime.GOOS == "linux" {
		openCmd := exec.Command("xdg-open", url)
		return openCmd.Run()
	}

	return browser.OpenURL(url)
}
