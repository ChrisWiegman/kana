package site

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/settings"

	"github.com/pkg/browser"
)

type Site struct {
	dockerClient           *docker.Client
	maxVerificationRetries int
	settings               *settings.Settings
	Named                  bool
}

type SiteInfo struct {
	Name, Path string
	Running    bool
}

const DefaultType = "site"

func Load(site *Site, kanaSettings *settings.Settings) {
	site.settings = kanaSettings
}

// EnsureDocker Ensures Docker is available for commands that need it.
func (s *Site) EnsureDocker(consoleOutput *console.Console) error {
	// Add a docker client to the site
	dockerClient, err := docker.New(consoleOutput, s.settings.Get("appDirectory"))
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

	code, checkURL, err := s.WPCli(checkCommand, false, consoleOutput)
	if err != nil || code != 0 {
		return fmt.Errorf("unable to determine SSL status")
	}

	parsedURL, err := url.Parse(strings.TrimSpace(checkURL))
	if err != nil {
		return err
	}

	if parsedURL.Scheme == "https" {
		localSettings["ssl"] = true
	}

	return s.settings.WriteLocalSettings(localSettings)
}

// GetSiteLink returns the link to the site.
func (s *Site) GetSiteLink() (string, error) {
	siteList, err := s.GetSiteList(false)
	if err != nil {
		return "", err
	}

	for i := range siteList {
		if siteList[i].Name == s.settings.Get("name") {
			return siteList[i].Path, nil
		}
	}

	return "", fmt.Errorf("site path not found")
}

// GetSiteList Returns a list of all Kana sites, their location and whether they're running.
func (s *Site) GetSiteList(checkRunningStatus bool) ([]SiteInfo, error) {
	sites := []SiteInfo{}

	sitesDir := filepath.Join(s.settings.Get("appDirectory"), "sites")

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

		content, err := os.ReadFile(filepath.Join(sitesDir, f.Name(), "link.json"))
		if err != nil {
			return sites, err
		}

		var jsonLink map[string]interface{}
		err = json.Unmarshal(content, &jsonLink)
		if err != nil {
			return sites, err
		}

		siteInfo.Name = f.Name()

		sitePath := fmt.Sprint(jsonLink["link"])

		if !strings.Contains(sitePath, sitesDir) {
			siteInfo.Path = sitePath
		}

		if checkRunningStatus {
			containers, err := s.dockerClient.ContainerList(f.Name())
			if err != nil {
				return sites, err
			}

			siteInfo.Running = len(containers) != 0
		}

		sites = append(sites, siteInfo)
	}

	return sites, nil
}

// IsSiteRunning Returns true if the site is up and running in Docker or false. Does not verify other errors.
func (s *Site) IsSiteRunning() bool {
	containers, _ := s.dockerClient.ContainerList(s.settings.Get("name"))

	return len(containers) != 0
}

// OpenSite Opens the current site in a browser if it is running.
func (s *Site) OpenSite(openDatabaseFlag, openMailpitFlag, openSiteFlag, openAdminFlag bool, consoleOutput *console.Console) error {
	openUrls := []string{}

	if openSiteFlag {
		openUrls = append(openUrls, s.settings.GetURL())
	}

	if openAdminFlag {
		openUrls = append(openUrls, s.settings.GetURL()+"/wp-admin/")
	}

	if openDatabaseFlag {
		isUsingSQLite, err := s.isUsingSQLite()
		if err != nil {
			return err
		}

		if isUsingSQLite {
			consoleOutput.Warn(fmt.Sprintf(
				"SQLite databases do not have a web interface and cannot be opened in TablePlus by URL. Open the database file, %s, directly using your database client of choice.", //nolint:lll
				filepath.Join(s.settings.Get("workingDirectory"), "wp-content", "database", ".ht.sqlite")))
			os.Exit(0)
		}

		databasePort := s.getDatabasePort()

		databaseURL := fmt.Sprintf(
			"mysql://wordpress:wordpress@127.0.0.1:%s/wordpress",
			databasePort)

		if s.settings.Get("databaseClient") == "phpmyadmin" {
			err := s.startPHPMyAdmin(consoleOutput)
			if err != nil {
				return err
			}

			databaseURL = fmt.Sprintf("%s://phpmyadmin-%s", s.settings.GetProtocol(), s.settings.GetDomain())
		}

		openUrls = append(openUrls, databaseURL)
	}

	if openMailpitFlag {
		if !s.isMailpitRunning() {
			err := s.startMailpit(consoleOutput)
			if err != nil {
				return err
			}
		}

		mailpitURL := fmt.Sprintf("%s://mailpit-%s", s.settings.GetProtocol(), s.settings.GetDomain())
		openUrls = append(openUrls, mailpitURL)
	}

	for _, openURL := range openUrls {
		var err error

		if strings.HasPrefix(openURL, "http") {
			err = s.verifySite(openURL)
			if err != nil {
				return err
			}
		}

		if runtime.GOOS == "linux" {
			openCmd := Command("xdg-open", openURL)
			return openCmd.Run()
		}

		err = browser.OpenURL(openURL)
		if err != nil {
			return err
		}
	}

	return nil
}

// StartSite Starts a site, including Traefik if needed.
func (s *Site) StartSite(consoleOutput *console.Console) error {
	// Let's start everything up
	consoleOutput.Printf("Starting development site: %s.\n", consoleOutput.Bold(consoleOutput.Green(s.settings.GetURL())))

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

	// Make sure we're using the correct database server
	err = s.maybeSetupSQLite()
	if err != nil {
		return err
	}

	// Install the Kana development plugin
	err = s.installKanaPlugin()
	if err != nil {
		return err
	}

	// Start Mailpit
	if s.settings.GetBool("mailpit") {
		err = s.startMailpit(consoleOutput)
		if err != nil {
			return err
		}
	}

	// Make sure the WordPress site is running
	err = s.verifySite(s.settings.GetURL())
	if err != nil {
		return err
	}

	// Setup WordPress
	err = s.installWordPress(consoleOutput)
	if err != nil {
		return err
	}

	// Verify the WordPress file permissions are correct
	err = s.resetWPFilePermissions()
	if err != nil {
		return err
	}

	// Maybe Remove the default plugins
	err = s.maybeRemoveDefaultPlugins()
	if err != nil {
		return err
	}

	// Install Xdebug if we need to
	if s.settings.GetBool("xdebug") {
		consoleOutput.Println("Installing and configuring Xdebug.")

		err = s.StartXdebug(consoleOutput)
		if err != nil {
			return err
		}
	}

	// Install any configuration plugins if needed
	err = s.installDefaultPlugins(consoleOutput)
	if err != nil {
		return err
	}

	// Activate the default theme if set
	err = s.activateTheme(consoleOutput)
	if err != nil {
		return err
	}

	// Activate the current project if asked
	err = s.activateProject(consoleOutput)
	if err != nil {
		return err
	}

	// Open the site in the user's browser
	return s.OpenSite(false, false, true, false, consoleOutput)
}

// StopSite Stops a full site, including Traefik if needed.
func (s *Site) StopSite() error {
	err := s.stopWordPress()
	if err != nil {
		return err
	}

	// If no other sites are running, also shut down the Traefik container
	return s.maybeStopTraefik()
}

// getDirectories Returns the correct appDir and databaseDir for the current site.
func (s *Site) getDirectories() (wordPressDirectory, databaseDir string, err error) {
	wordPressDirectory, err = s.getWordPressDirectory()
	if err != nil {
		return "", "", err
	}

	databaseDir, err = s.getDatabaseDirectory()
	if err != nil {
		return "", "", err
	}

	return wordPressDirectory, databaseDir, err
}

// getRunningConfig gets various options that were used to start the site.
func (s *Site) getRunningConfig(withPlugins bool, consoleOutput *console.Console) (map[string]interface{}, error) {
	localSettings := s.settings.GetAll("local")

	// We need container details to see if the mailpit container is running
	localSettings["mailpit"] = s.isMailpitRunning()

	output, err := s.WordPress("pecl list | grep xdebug", false, false)
	if err != nil {
		return localSettings, err
	}

	if strings.Contains(output.StdOut, "xdebug") {
		localSettings["xdebug"] = true
	}

	output, err = s.WordPress("echo $WORDPRESS_DEBUG", false, false)
	if err != nil {
		return localSettings, err
	}

	if strings.Contains(output.StdOut, "1") {
		localSettings["wpdebug"] = true
	}

	output, err = s.WordPress("echo $SCRIPT_DEBUG", false, false)
	if err != nil {
		return localSettings, err
	}

	if strings.Contains(output.StdOut, "1") {
		localSettings["scriptDebug"] = true
	}

	output, err = s.WordPress("echo $KANA_SQLITE", false, false)
	if err != nil {
		return localSettings, err
	}

	if strings.Contains(output.StdOut, "true") {
		localSettings["database"] = "sqlite"
	}

	mounts := s.dockerClient.ContainerGetMounts(fmt.Sprintf("kana-%s-wordpress", s.settings.Get("name")))

	if len(mounts) == 1 {
		localSettings["type"] = DefaultType
	}

	for _, mount := range mounts {
		if strings.Contains(mount.Destination, "/var/www/html/wp-content/plugins/") {
			localSettings["type"] = "plugin"
		}

		if strings.Contains(mount.Destination, "/var/www/html/wp-content/themes/") {
			localSettings["type"] = "theme" //nolint:goconst
		}
	}

	// Don't get plugins if we don't need them
	if withPlugins {
		plugins, hasDefaultPlugins, err := s.getInstalledWordPressPlugins(consoleOutput)
		if err != nil {
			return localSettings, err
		}

		localSettings["removeDefaultPlugins"] = !hasDefaultPlugins
		localSettings["plugins"] = plugins
	}

	return localSettings, nil
}

// maybeRemoveDefaultPlugins Removes the default plugins if the setting is set.
func (s *Site) maybeRemoveDefaultPlugins() error {
	if !s.settings.GetBool("RemoveDefaultPlugins") {
		return nil
	}

	wordPressDirectory, err := s.getWordPressDirectory()
	if err != nil {
		return err
	}

	defaultPlugins := []string{
		"hello.php",
		"akismet"}

	for _, plugin := range defaultPlugins {
		pluginPath := filepath.Join(wordPressDirectory, "wp-content", "plugins", plugin)
		err = os.RemoveAll(pluginPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// startContainer Starts a given container configuration.
func (s *Site) startContainer(container *docker.ContainerConfig, randomPorts, localUser bool, consoleOutput *console.Console) error {
	err := s.dockerClient.EnsureImage(container.Image, s.settings.Get("appDirectory"), s.settings.GetInt("updateInterval"), consoleOutput)
	if err != nil {
		err = s.handleImageError(container, err)
		if err != nil {
			return err
		}
	}
	_, err = s.dockerClient.ContainerRun(container, randomPorts, localUser)

	return err
}

// verifySite verifies if a site is up and running without error.
func (s *Site) verifySite(siteURL string) error {
	// Setup other options generated from config items
	rootCert := filepath.Join(s.settings.Get("appDirectory"), "certs", s.settings.Get("rootCert"))

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

		if tries == s.maxVerificationRetries {
			return errors.New("timeout reached. unable to open site")
		}

		tries++
		time.Sleep(time.Second)
	}

	return nil
}
