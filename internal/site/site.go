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

	"github.com/ChrisWiegman/kana-cli/internal/config"
	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/docker"

	"github.com/pkg/browser"
)

type Site struct {
	dockerClient *docker.DockerClient
	config       *config.Config
}

// NewSite creates a new site object
func NewSite(kanaConfig *config.Config) (*Site, error) {

	site := new(Site)

	// Add a docker client to the site
	dockerClient, err := docker.NewController()
	if err != nil {
		return site, err
	}

	site.dockerClient = dockerClient

	site.config = kanaConfig

	return site, nil
}

// GetURL returns the appropriate URL for the site
func (s *Site) GetURL(insecure bool) string {

	if insecure {
		return s.config.Site.URL
	}

	return s.config.Site.SecureURL
}

// VerifySite verifies if a site is up and running without error
func (s *Site) VerifySite() (bool, error) {

	// Setup other options generated from config items
	rootCert := path.Join(s.config.Directories.App, "certs", s.config.App.RootCert)

	caCert, err := os.ReadFile(rootCert)
	if err != nil {
		return false, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(s.config.Site.SecureURL)
	if err != nil {
		return false, err
	}

	tries := 0

	for resp.StatusCode != 200 {

		resp, err = client.Get(s.config.Site.SecureURL)
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

	openURL(s.config.Site.SecureURL)

	return nil
}

// InstallXdebug installs xdebug in the site's PHP container
func (s *Site) InstallXdebug() (bool, error) {

	if !s.config.Site.Xdebug {
		return false, nil
	}

	console.Println("Installing Xdebug...")

	commands := []string{
		"pecl list | grep xdebug",
		"pecl install xdebug",
		"docker-php-ext-enable xdebug",
		"echo 'xdebug.start_with_request=yes' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.mode=debug' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.client_host=host.docker.internal' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.discover_client_host=on' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.start_with_request=trigger' >> /usr/local/etc/php/php.ini",
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

	container := fmt.Sprintf("kana_%s_wordpress", s.config.Site.SiteName)

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

// IsLocalSite Determines if a site is a "local" site (started with the "local" flag) so that other commands can work as needed.
func (s *Site) IsLocalSite() bool {

	// If the site is already running, try to make this easier
	if s.IsSiteRunning() {
		runningConfig := s.GetRunningConfig()
		if runningConfig.Local {
			return true
		}
	}

	// First check the app site folders. If we've created the site (has a DB) without an "app" folder we can assume it was local last time.
	hasNonLocalAppFolder := true
	hasDatabaseFolder := true

	if _, err := os.Stat(path.Join(s.config.Directories.Site, "app")); os.IsNotExist(err) {
		hasNonLocalAppFolder = false
	}

	if _, err := os.Stat(path.Join(s.config.Directories.Site, "database")); os.IsNotExist(err) {
		hasDatabaseFolder = false
	}

	if hasDatabaseFolder && !hasNonLocalAppFolder {
		return true
	}

	// Return the flag for all other conditions
	return s.config.Site.Local
}

// GetRunningConfig gets various options that were used to start the site
func (s *Site) GetRunningConfig() CurrentConfig {

	currentConfig := CurrentConfig{
		Type:   "site",
		Local:  false,
		Xdebug: false,
	}

	output, _ := s.runCli("pecl list | grep xdebug", false)
	if strings.Contains(output.StdOut, "xdebug") {
		currentConfig.Xdebug = true
	}

	mounts := s.dockerClient.ContainerGetMounts(fmt.Sprintf("kana_%s_wordpress", s.config.Site.SiteName))

	if len(mounts) == 1 {
		currentConfig.Type = "site"
	}

	for _, mount := range mounts {

		if mount.Source == path.Join(s.config.Directories.Working, "wordpress") {
			currentConfig.Local = true
		}

		if strings.Contains(mount.Destination, "/var/www/html/wp-content/plugins/") {
			currentConfig.Type = "plugin"
		}

		if strings.Contains(mount.Destination, "/var/www/html/wp-content/themes/") {
			currentConfig.Type = "theme"
		}
	}

	return currentConfig
}

func (s *Site) ExportSiteConfig() error {

	config := s.GetRunningConfig()
	plugins, err := s.GetInstalledWordPressPlugins()
	if err != nil {
		return err
	}

	s.config.Site.Viper.Set("local", config.Local)
	s.config.Site.Viper.Set("type", config.Type)
	s.config.Site.Viper.Set("xdebug", config.Xdebug)
	s.config.Site.Viper.Set("plugins", plugins)

	if _, err = os.Stat(path.Join(s.config.Directories.Working, ".kana.json")); os.IsNotExist(err) {
		return s.config.Site.Viper.SafeWriteConfig()
	}

	return s.config.Site.Viper.WriteConfig()
}
