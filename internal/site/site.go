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
	"github.com/spf13/cobra"
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

// NewSite creates a new site object
func NewSite(staticConfig appConfig.StaticConfig, dynamicConfig *viper.Viper) (*Site, error) {

	site := new(Site)

	// Add a docker client to the site
	dockerClient, err := docker.NewController()
	if err != nil {
		return site, err
	}

	site.dockerClient = dockerClient

	// Setup all config items (static, dynamic and site options)
	site.StaticConfig = staticConfig
	site.DynamicConfig = dynamicConfig
	site.SiteConfig, err = getSiteConfig(staticConfig, dynamicConfig)
	if err != nil {
		return site, err
	}

	// Setup other options generated from config items
	site.rootCert = path.Join(staticConfig.AppDirectory, "certs", staticConfig.RootCert)
	site.siteDomain = fmt.Sprintf("%s.%s", staticConfig.SiteName, staticConfig.AppDomain)
	site.secureURL = fmt.Sprintf("https://%s/", site.siteDomain)
	site.url = fmt.Sprintf("http://%s/", site.siteDomain)

	return site, nil
}

// ProcessNameFlag Processes the name flag on the site resetting all appropriate site variables
func (s *Site) ProcessNameFlag(cmd *cobra.Command, siteName string) {

	if !cmd.Flags().Lookup("name").Changed {
		return
	}

	// Check that we're not using invalid start flags for the start command
	if cmd.Use == "start" {
		if cmd.Flags().Lookup("plugin").Changed || cmd.Flags().Lookup("theme").Changed || cmd.Flags().Lookup("local").Changed {
			return
		}
	}

	fmt.Println(s.StaticConfig.SiteName)
	fmt.Println(siteName)

	s.StaticConfig.SiteName = siteName
	s.StaticConfig.SiteDirectory = (path.Join(s.StaticConfig.AppDirectory, "sites", siteName))

	s.siteDomain = fmt.Sprintf("%s.%s", s.StaticConfig.SiteName, s.StaticConfig.AppDomain)
	s.secureURL = fmt.Sprintf("https://%s/", s.siteDomain)
	s.url = fmt.Sprintf("http://%s/", s.siteDomain)
}

// GetURL returns the appropriate URL for the site
func (s *Site) GetURL(insecure bool) string {

	if insecure {
		return s.url
	}

	return s.secureURL
}

// VerifySite verifies if a site is up and running without error
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

// OpenSite Opens the current site in a browser if it is running correctly
func (s *Site) OpenSite() error {

	_, err := s.VerifySite()
	if err != nil {
		return err
	}

	openURL(s.secureURL)

	return nil
}

// InstallXdebug installs xdebug in the site's PHP container
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

		// Verify that the command ran correctly
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

// openURL opens the URL in the user's default browser based on which OS they're using
func openURL(url string) error {

	if runtime.GOOS == "linux" {
		openCmd := exec.Command("xdg-open", url)
		return openCmd.Run()
	}

	return browser.OpenURL(url)
}
