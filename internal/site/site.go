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
	Config       *config.Config
}

// NewSite creates a new site object
func NewSite() (*Site, error) {

	site := new(Site)

	// Add a docker client to the site
	dockerClient, err := docker.NewController()
	if err != nil {
		return site, err
	}

	site.dockerClient = dockerClient

	err = site.loadConfig()

	return site, err
}

// ExportSiteSConfig Saves the current running config to a file.
func (s *Site) ExportSiteConfig() error {

	config := s.getRunningConfig()
	plugins, err := s.getInstalledWordPressPlugins()
	if err != nil {
		return err
	}

	s.Config.Local.Viper.Set("local", config.Local)
	s.Config.Local.Viper.Set("type", config.Type)
	s.Config.Local.Viper.Set("xdebug", config.Xdebug)
	s.Config.Local.Viper.Set("plugins", plugins)

	if _, err = os.Stat(path.Join(s.Config.Directories.Working, ".kana.json")); os.IsNotExist(err) {
		return s.Config.Local.Viper.SafeWriteConfig()
	}

	return s.Config.Local.Viper.WriteConfig()
}

// IsSiteRunning Returns true if the site is up and running in Docker or false. Does not verify other errors
func (s *Site) IsSiteRunning() bool {

	containers, _ := s.dockerClient.ListContainers(s.Config.Local.Name)

	return len(containers) != 0
}

// OpenSite Opens the current site in a browser if it is running
func (s *Site) OpenSite() error {

	_, err := s.verifySite()
	if err != nil {
		return err
	}

	if runtime.GOOS == "linux" {
		openCmd := exec.Command("xdg-open", s.Config.Local.SecureURL)
		return openCmd.Run()
	}

	return browser.OpenURL(s.Config.Local.SecureURL)
}

// StartSite Starts a site, including Traefik if needed
func (s *Site) StartSite() error {

	// Let's start everything up
	fmt.Printf("Starting development site: %s\n", s.getSiteURL(false))

	// Start Traefik if we need it
	err := s.startTraefik()
	if err != nil {
		return err
	}

	// Start WordPress
	err = s.startWordPress()
	if err != nil {
		return err
	}

	// Make sure the WordPress site is running
	_, err = s.verifySite()
	if err != nil {
		return err
	}

	// Setup WordPress
	err = s.installWordPress()
	if err != nil {
		return err
	}

	// Install Xdebug if we need to
	_, err = s.installXdebug()
	if err != nil {
		return err
	}

	// Install any configuration plugins if needed
	err = s.installDefaultPlugins()
	if err != nil {
		return err
	}

	// Open the site in the user's browser
	return s.OpenSite()
}

// StopSite Stops a full site, including Traefik if needed
func (s *Site) StopSite() error {

	err := s.stopWordPress()
	if err != nil {
		return err
	}

	// If no other sites are running, also shut down the Traefik container
	return s.maybeStopTraefik()
}

// getLocalAppDir Gets the absolute path to WordPress if the local flag or option has been set
func getLocalAppDir() (string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	localAppDir := path.Join(cwd, "wordpress")

	err = os.MkdirAll(localAppDir, 0750)
	if err != nil {
		return "", err
	}

	return localAppDir, nil
}

// getRunningConfig gets various options that were used to start the site
func (s *Site) getRunningConfig() CurrentConfig {

	currentConfig := CurrentConfig{
		Type:   "site",
		Local:  false,
		Xdebug: false,
	}

	output, _ := s.runCli("pecl list | grep xdebug", false)
	if strings.Contains(output.StdOut, "xdebug") {
		currentConfig.Xdebug = true
	}

	mounts := s.dockerClient.ContainerGetMounts(fmt.Sprintf("kana_%s_wordpress", s.Config.Local.Name))

	if len(mounts) == 1 {
		currentConfig.Type = "site"
	}

	for _, mount := range mounts {

		if mount.Source == path.Join(s.Config.Directories.Working, "wordpress") {
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

// getSiteURL returns the appropriate URL for the site
func (s *Site) getSiteURL(insecure bool) string {

	if insecure {
		return s.Config.Local.URL
	}

	return s.Config.Local.SecureURL
}

// installXdebug installs xdebug in the site's PHP container
func (s *Site) installXdebug() (bool, error) {

	if !s.Config.Local.Xdebug {
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

// isLocalSite Determines if a site is a "local" site (started with the "local" flag) so that other commands can work as needed.
func (s *Site) isLocalSite() bool {

	// If the site is already running, try to make this easier
	if s.IsSiteRunning() {
		runningConfig := s.getRunningConfig()
		if runningConfig.Local {
			return true
		}
	}

	// First check the app site folders. If we've created the site (has a DB) without an "app" folder we can assume it was local last time.
	hasNonLocalAppFolder := true
	hasDatabaseFolder := true

	if _, err := os.Stat(path.Join(s.Config.Directories.Site, "app")); os.IsNotExist(err) {
		hasNonLocalAppFolder = false
	}

	if _, err := os.Stat(path.Join(s.Config.Directories.Site, "database")); os.IsNotExist(err) {
		hasDatabaseFolder = false
	}

	if hasDatabaseFolder && !hasNonLocalAppFolder {
		return true
	}

	// Return the flag for all other conditions
	return s.Config.Local.Local
}

func (s *Site) loadConfig() error {

	var err error

	s.Config, err = config.NewConfig()
	return err
}

// runCli Runs an arbitrary CLI command against the site's WordPress container
func (s *Site) runCli(command string, restart bool) (docker.ExecResult, error) {

	container := fmt.Sprintf("kana_%s_wordpress", s.Config.Local.Name)

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

// verifySite verifies if a site is up and running without error
func (s *Site) verifySite() (bool, error) {

	// Setup other options generated from config items
	rootCert := path.Join(s.Config.Directories.App, "certs", s.Config.Global.RootCert)

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

	resp, err := client.Get(s.Config.Local.SecureURL)
	if err != nil {
		return false, err
	}

	tries := 0

	for resp.StatusCode != 200 {

		resp, err = client.Get(s.Config.Local.SecureURL)
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
