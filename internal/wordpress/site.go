package wordpress

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

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/pkg/browser"
)

type KanaSite struct {
	controller *docker.Controller
	rootCert   string
	siteDomain string
	secureURL  string
	url        string
}

func NewSite(controller *docker.Controller) *KanaSite {

	site := new(KanaSite)

	site.controller = controller
	site.rootCert = path.Join(controller.Config.AppDirectory, "certs", controller.Config.RootCert)
	site.siteDomain = fmt.Sprintf("%s.%s", controller.Config.SiteDirectory, controller.Config.AppDomain)
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

func (s *KanaSite) VerifySite() (bool, error) {

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

func (s *KanaSite) OpenSite() error {

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
