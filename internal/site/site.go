package site

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/ChrisWiegman/kana/internal/appConfig"
	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

type Site struct {
	dockerClient  *docker.DockerClient
	StaticConfig  appConfig.StaticConfig
	DynamicConfig *viper.Viper
	SiteConfig    *viper.Viper
	rootCert      string
	siteDomain    string
	secureURL     string
	url           string
}

func NewSite(staticConfig appConfig.StaticConfig, dynamicConfig *viper.Viper) (*Site, error) {

	site := new(Site)

	dockerClient, err := docker.NewController()
	if err != nil {
		return site, err
	}

	site.StaticConfig = staticConfig
	site.SiteConfig, err = getSiteConfig(staticConfig, dynamicConfig)
	if err != nil {
		return site, err
	}

	site.dockerClient = dockerClient
	site.rootCert = path.Join(staticConfig.AppDirectory, "certs", staticConfig.RootCert)
	site.siteDomain = fmt.Sprintf("%s.%s", staticConfig.SiteName, staticConfig.AppDomain)
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

	caCert, err := os.ReadFile(s.rootCert)
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

func (s *Site) InstallXdebug() (bool, error) {

	if !s.SiteConfig.GetBool("xdebug") {
		return false, nil
	}

	fmt.Println("Installing Xdebug...")

	commands := []string{
		"pecl list | grep xdebug",
		"pecl install xdebug",
		"docker-php-ext-enable xdebug",
		"echo 'xdebug.start_with_request=yes' >> /usr/local/etc/php/php.ini",
		"xdebug.mode=debug' >> /usr/local/etc/php/php.ini",
	}

	for i, command := range commands {

		restart := false

		if i+1 == len(commands) {
			restart = true
		}

		output, err := s.runCli(command, restart)
		if err != nil {
			return false, err
		}

		if i == 0 && strings.Contains(output.StdOut, "xdebug") {
			return false, nil
		}
	}

	return true, nil

}

// runCli Runs an arbitrary CLI command against the site's WordPress container
func (s *Site) runCli(command string, restart bool) (docker.ExecResult, error) {

	container := fmt.Sprintf("kana_%s_wordpress", s.StaticConfig.SiteName)

	output, err := s.dockerClient.ContainerExec(container, []string{command})
	if err != nil {
		return docker.ExecResult{}, err
	}

	if restart {
		_, err = s.dockerClient.ContainerRestart(container)
		return output, err
	}

	return output, nil

}

func openURL(url string) error {

	if runtime.GOOS == "linux" {
		openCmd := exec.Command("xdg-open", url)
		return openCmd.Run()
	}

	return browser.OpenURL(url)
}
