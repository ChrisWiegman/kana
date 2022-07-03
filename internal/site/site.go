package site

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path"
	"runtime"
	"time"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/pkg/browser"
)

type Site struct {
	dockerClient *docker.DockerClient
	appConfig    config.AppConfig
	rootCert     string
	siteDomain   string
	secureURL    string
	url          string
}

func NewSite(appConfig config.AppConfig) (*Site, error) {

	site := new(Site)

	dockerClient, err := docker.NewController()
	if err != nil {
		return site, err
	}

	site.appConfig = appConfig
	site.dockerClient = dockerClient
	site.rootCert = path.Join(appConfig.AppDirectory, "certs", appConfig.RootCert)
	site.siteDomain = fmt.Sprintf("%s.%s", appConfig.SiteDirectory, appConfig.AppDomain)
	site.secureURL = fmt.Sprintf("https://%s/", site.siteDomain)
	site.url = fmt.Sprintf("http://%s/", site.siteDomain)

	return site, nil
}

func (s *Site) GetURL(insecure bool) string {

	if insecure {
		return s.url
	}

	return s.secureURL

}

func (s *Site) VerifySite() (bool, error) {

	caCert, err := ioutil.ReadFile(s.rootCert)
	if err != nil {
		return false, err
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
		return false, err
	}

	tries := 0

	for resp.StatusCode != 200 {

		resp, err = client.Get(s.secureURL)
		if err != nil {
			return false, err
		}

		if resp.StatusCode == 200 {
			break
		}

		if tries == 30 {
			return false, fmt.Errorf("timeout reached. unable to open site")
		}

		tries++
		time.Sleep(1 * time.Second)

	}

	return true, nil

}

func (s *Site) OpenSite() error {

	_, err := s.VerifySite()
	if err != nil {
		return err
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
