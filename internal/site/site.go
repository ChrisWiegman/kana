package site

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/docker"
	"github.com/ChrisWiegman/kana-cli/internal/settings"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

type Site struct {
	dockerClient *docker.DockerClient
	Settings     *settings.Settings
}

type SiteInfo struct {
	Name, Path string
	Running    bool
}

var maxVerificationRetries = 30

var execCommand = exec.Command

// EnsureDocker Ensures Docker is available for commands that need it.
func (s *Site) EnsureDocker(consoleOutput *console.Console) error {
	// Add a docker client to the site
	dockerClient, err := docker.NewDockerClient(consoleOutput, s.Settings.AppDirectory)
	if err != nil {
		return err
	}

	s.dockerClient = dockerClient
	return nil
}

// ExportSiteSConfig Saves the current running config to a file.
func (s *Site) ExportSiteConfig(consoleOutput *console.Console) error {
	localSettings, err := s.getRunningConfig(true, consoleOutput)
	if err != nil {
		return err
	}

	checkCommand := []string{
		"option",
		"get",
		"siteurl",
	}

	code, checkURL, err := s.RunWPCli(checkCommand, consoleOutput)
	if err != nil || code != 0 {
		return fmt.Errorf("unable to determine SSL status")
	}

	parsedURL, err := url.Parse(strings.TrimSpace(checkURL))
	if err != nil {
		return err
	}

	if parsedURL.Scheme == "https" {
		localSettings.SSL = true
	}

	return s.Settings.WriteLocalSettings(localSettings)
}

// GetSiteList Returns a list of all Kana sites, their location and whether they're running
func (s *Site) GetSiteList(appDir string, consoleOutput *console.Console) ([]SiteInfo, error) {
	var sites []SiteInfo
	sitesDir := path.Join(appDir, "sites")

	_, err := os.Stat(sitesDir)
	if os.IsNotExist(err) {
		return sites, nil
	}

	appSites, err := os.ReadDir(sitesDir)
	if err != nil {
		return sites, err
	}

	for _, f := range appSites {
		var siteInfo SiteInfo

		content, err := os.ReadFile(path.Join(sitesDir, f.Name(), "link.json"))
		if err != nil {
			return sites, err
		}

		var jsonLink map[string]interface{}
		err = json.Unmarshal(content, &jsonLink)
		if err != nil {
			return sites, err
		}

		containers, err := s.dockerClient.ContainerList(f.Name())
		if err != nil {
			return sites, err
		}

		sitePath := fmt.Sprint(jsonLink["link"])

		if !strings.Contains(sitePath, sitesDir) {
			siteInfo.Path = sitePath
		}

		siteInfo.Name = f.Name()
		siteInfo.Running = len(containers) != 0

		sites = append(sites, siteInfo)
	}

	return sites, nil
}

// IsSiteRunning Returns true if the site is up and running in Docker or false. Does not verify other errors
func (s *Site) IsSiteRunning() bool {
	containers, _ := s.dockerClient.ContainerList(s.Settings.Name)

	return len(containers) != 0
}

func (s *Site) LoadSite(cmd *cobra.Command, commandsRequiringSite []string, startFlags settings.StartFlags, flagVerbose bool) error {
	var err error

	s.Settings, err = settings.NewSettings()
	if err != nil {
		return err
	}

	// Load app-wide settings
	err = s.Settings.LoadGlobalSettings()
	if err != nil {
		return err
	}

	// Load settings specific to the site
	isSite, err := s.Settings.LoadLocalSettings(cmd)
	if err != nil {
		return err
	}

	// Fail now if we have a command that requires a completed site and we haven't started it before
	if !isSite && arrayContains(commandsRequiringSite, cmd.Use) {
		return fmt.Errorf("the current site you are trying to work with does not exist. Use `kana start` to initialize")
	}

	// Process the "start" command flags
	if cmd.Use == "start" {
		// A site shouldn't be both a plugin and a theme so this reports an error if that is the case.
		if startFlags.IsPlugin && startFlags.IsTheme {
			return fmt.Errorf("you have set both the plugin and theme flags. Please choose only one option")
		}

		s.Settings.ProcessStartFlags(cmd, startFlags)
	}

	return nil
}

// OpenSite Opens the current site in a browser if it is running
func (s *Site) OpenSite(openPhpMyAdminFlag, openMailpitFlag, openSiteFlag bool, consoleOutput *console.Console) error {
	openUrls := []string{}

	if openSiteFlag {
		openUrls = append(openUrls, s.Settings.URL)
	}

	if openPhpMyAdminFlag {
		err := s.startPHPMyAdmin(consoleOutput)
		if err != nil {
			return err
		}

		phpmyAdminURL := fmt.Sprintf("%s://phpmyadmin-%s", s.Settings.Protocol, s.Settings.SiteDomain)
		openUrls = append(openUrls, phpmyAdminURL)
	}

	if openMailpitFlag {
		if !s.isMailpitRunning() {
			err := s.startMailpit(consoleOutput)
			if err != nil {
				return err
			}
		}

		mailpitURL := fmt.Sprintf("%s://mailpit-%s", s.Settings.Protocol, s.Settings.SiteDomain)
		openUrls = append(openUrls, mailpitURL)
	}

	for _, openURL := range openUrls {
		err := s.verifySite(openURL)
		if err != nil {
			return err
		}

		if runtime.GOOS == "linux" {
			openCmd := execCommand("xdg-open", openURL)
			return openCmd.Run()
		}

		err = browser.OpenURL(openURL)
		if err != nil {
			return err
		}
	}

	return nil
}

// StartSite Starts a site, including Traefik if needed
func (s *Site) StartSite(consoleOutput *console.Console) error {
	// Let's start everything up
	consoleOutput.Printf("Starting development site: %s.\n", consoleOutput.Bold(consoleOutput.Green(s.Settings.URL)))

	// Start Traefik if we need it
	err := s.startTraefik(consoleOutput)
	if err != nil {
		return err
	}

	// Start WordPress
	err = s.startWordPress(consoleOutput)
	if err != nil {
		return err
	}

	// Start Mailpit
	if s.Settings.Mailpit {
		err = s.startMailpit(consoleOutput)
		if err != nil {
			return err
		}
	}

	// Make sure the WordPress site is running
	err = s.verifySite(s.Settings.URL)
	if err != nil {
		return err
	}

	// Setup WordPress
	err = s.installWordPress(consoleOutput)
	if err != nil {
		return err
	}

	// Install Xdebug if we need to
	_, err = s.installXdebug(consoleOutput)
	if err != nil {
		return err
	}

	// Install the Kana development plugin
	err = s.installKanaPlugin(consoleOutput)
	if err != nil {
		return err
	}

	// Install any configuration plugins if needed
	err = s.installDefaultPlugins(consoleOutput)
	if err != nil {
		return err
	}

	// Open the site in the user's browser
	return s.OpenSite(false, false, true, consoleOutput)
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

// checkStatusCode returns true on 200 or false
func checkStatusCode(checkURL string) (bool, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, checkURL, http.NoBody)
	if err != nil {
		return false, err
	}

	// Ignore SSL check as we're using our self-signed cert for development
	clientTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
	}

	client := &http.Client{Transport: clientTransport}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

// getLocalAppDir Gets the absolute path to WordPress if the local flag or option has been set
func (s *Site) getLocalAppDir() (string, error) {
	localAppDir := path.Join(s.Settings.WorkingDirectory, "wordpress")

	err := os.MkdirAll(localAppDir, os.FileMode(defaultDirPermissions))
	if err != nil {
		return "", err
	}

	return localAppDir, nil
}

// getRunningConfig gets various options that were used to start the site
func (s *Site) getRunningConfig(withPlugins bool, consoleOutput *console.Console) (settings.LocalSettings, error) {
	localSettings := settings.LocalSettings{
		Type:    "site",
		Local:   false,
		Xdebug:  false,
		SSL:     false,
		Mailpit: false,
		WPDebug: false,
	}

	// We need container details to see if the mailpit container is running
	localSettings.Mailpit = s.isMailpitRunning()

	output, err := s.runCli("pecl list | grep xdebug", false)
	if err != nil {
		return localSettings, err
	}

	if strings.Contains(output.StdOut, "xdebug") {
		localSettings.Xdebug = true
	}

	output, err = s.runCli("echo $WORDPRESS_DEBUG", false)
	if err != nil {
		return localSettings, err
	}

	if strings.Contains(output.StdOut, "1") {
		localSettings.WPDebug = true
	}

	mounts := s.dockerClient.ContainerGetMounts(fmt.Sprintf("kana-%s-wordpress", s.Settings.Name))

	if len(mounts) == 1 {
		localSettings.Type = "site"
	}

	for _, mount := range mounts {
		if mount.Source == path.Join(s.Settings.WorkingDirectory, "wordpress") {
			localSettings.Local = true
		}

		if strings.Contains(mount.Destination, "/var/www/html/wp-content/plugins/") {
			localSettings.Type = "plugin" //nolint:goconst
		}

		if strings.Contains(mount.Destination, "/var/www/html/wp-content/themes/") {
			localSettings.Type = "theme" //nolint:goconst
		}
	}

	// Don't get plugins if we don't need them
	if withPlugins {
		plugins, err := s.getInstalledWordPressPlugins(consoleOutput)
		if err != nil {
			return localSettings, err
		}

		localSettings.Plugins = plugins
	}

	return localSettings, nil
}

func (s *Site) isMailpitRunning() bool {
	// We need container details to see if the mailpit container is running
	containers, err := s.dockerClient.ContainerList(s.Settings.Name)
	if err != nil {
		return false
	}

	for i := range containers {
		if containers[i].Image == "axllent/mailpit" {
			return true
		}
	}

	return false
}

// installXdebug installs xdebug in the site's PHP container
func (s *Site) installXdebug(consoleOutput *console.Console) (bool, error) {
	if !s.Settings.Xdebug {
		return false, nil
	}

	consoleOutput.Println("Installing and configuring Xdebug.")

	commands := []string{
		"pecl list | grep xdebug",
		"pecl install xdebug",
		"docker-php-ext-enable xdebug",
		"echo 'xdebug.start_with_request=yes' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.mode=debug,develop,trace' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.client_host=host.docker.internal' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.discover_client_host=on' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.start_with_request=trigger' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.show_local_vars=1' >> /usr/local/etc/php/php.ini",
		"echo 'html_errors = On' >> /usr/local/etc/php/conf.d/z-custom.ini", // Ensure custom overrides happen
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
func (s *Site) isLocalSite(consoleOutput *console.Console) bool {
	// If the site is already running, try to make this easier
	if s.IsSiteRunning() {
		runningConfig, _ := s.getRunningConfig(false, consoleOutput)
		if runningConfig.Local {
			return true
		}
	}

	// First check the app site folders. If we've created the site (has a DB) without an "app" folder we can assume it was local last time.
	hasNonLocalAppFolder := true
	hasDatabaseFolder := true

	if _, err := os.Stat(path.Join(s.Settings.SiteDirectory, "app")); os.IsNotExist(err) {
		hasNonLocalAppFolder = false
	}

	if _, err := os.Stat(path.Join(s.Settings.SiteDirectory, "database")); os.IsNotExist(err) {
		hasDatabaseFolder = false
	}

	if hasDatabaseFolder && !hasNonLocalAppFolder {
		return true
	}

	// Return the flag for all other conditions
	return s.Settings.Local
}

// runCli Runs an arbitrary CLI command against the site's WordPress container
func (s *Site) runCli(command string, restart bool) (docker.ExecResult, error) {
	container := fmt.Sprintf("kana-%s-wordpress", s.Settings.Name)

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
func (s *Site) verifySite(siteURL string) error {
	// Setup other options generated from config items
	rootCert := path.Join(s.Settings.AppDirectory, "certs", s.Settings.RootCert)

	caCert, err := os.ReadFile(rootCert)
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	siteOK, err := checkStatusCode(siteURL)
	if err != nil {
		return err
	}

	tries := 0

	for !siteOK {
		siteOK, err = checkStatusCode(siteURL)
		if err != nil {
			return err
		}

		if siteOK {
			break
		}

		if tries == maxVerificationRetries {
			return errors.New("timeout reached. unable to open site")
		}

		tries++
		time.Sleep(time.Second)
	}

	return nil
}
