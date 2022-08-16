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

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SiteConfig struct {
	PHPVersion string
	Xdebug     bool
	Local      bool
	Type       string
}

type StartFlags struct {
	Xdebug   bool
	Local    bool
	IsTheme  bool
	IsPlugin bool
}

type Site struct {
	dockerClient *docker.DockerClient
	appConfig    config.AppConfig
	siteConfig   SiteConfig
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
	site.siteConfig, err = getSiteConfig(appConfig)
	if err != nil {
		return site, err
	}

	site.dockerClient = dockerClient
	site.rootCert = path.Join(appConfig.AppDirectory, "certs", appConfig.RootCert)
	site.siteDomain = fmt.Sprintf("%s.%s", appConfig.SiteName, appConfig.AppDomain)
	site.secureURL = fmt.Sprintf("https://%s/", site.siteDomain)
	site.url = fmt.Sprintf("http://%s/", site.siteDomain)

	return site, nil
}

func getSiteConfig(appConfig config.AppConfig) (SiteConfig, error) {

	viperConfig := viper.New()

	viperConfig.SetDefault("php", appConfig.DefaultPHPVersion)
	viperConfig.SetDefault("type", appConfig.SiteType)
	viperConfig.SetDefault("local", appConfig.SiteLocal)
	viperConfig.SetDefault("xdebug", appConfig.SiteXdebug)

	viperConfig.SetConfigName(".kana")
	viperConfig.SetConfigType("json")
	viperConfig.AddConfigPath(appConfig.WorkingDirectory)

	err := viperConfig.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return SiteConfig{}, err
		}
	}

	siteConfig := SiteConfig{
		PHPVersion: viperConfig.GetString("php"),
		Type:       viperConfig.GetString("type"),
		Local:      viperConfig.GetBool("local"),
		Xdebug:     viperConfig.GetBool("xdebug"),
	}

	return siteConfig, nil

}

func (s *Site) AddStartFlags(cmd *cobra.Command, flags StartFlags) {

	if cmd.Flags().Lookup("local").Changed {
		s.siteConfig.Local = flags.Local
	}

	if cmd.Flags().Lookup("xdebug").Changed {
		s.siteConfig.Xdebug = flags.Xdebug
	}

	if cmd.Flags().Lookup("plugin").Changed && flags.IsPlugin {
		s.siteConfig.Type = "plugin"
	}

	if cmd.Flags().Lookup("theme").Changed && flags.IsTheme {
		s.siteConfig.Type = "theme"
	}
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

	if !s.siteConfig.Xdebug {
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

	container := fmt.Sprintf("kana_%s_wordpress", s.appConfig.SiteName)

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
