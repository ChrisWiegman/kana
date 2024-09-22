package site

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ChrisWiegman/kana-wordpress/internal/console"
	"github.com/ChrisWiegman/kana-wordpress/internal/docker"
	"github.com/ChrisWiegman/kana-wordpress/internal/settings"

	"github.com/docker/docker/api/types/mount"
)

type PluginInfo struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Update  string `json:"update"`
	Version string `json:"version"`
}

var defaultDirPermissions = 0750

func (s *Site) getWordPressDirectory() (wordPressDirectory string, err error) {
	wordPressDirectory = filepath.Join(s.settings.Get("siteDirectory"), "wordpress")

	siteLink, err := s.GetSiteLink()
	if err != nil {
		return "", err
	}

	if !s.settings.GetBool("isNamed") || siteLink != "" {
		wordPressDirectory = s.settings.Get("workingDirectory")

		if s.settings.Get("type") != DefaultType {
			wordPressDirectory = filepath.Join(s.settings.Get("workingDirectory"), "wordpress")
		}
	}

	err = os.MkdirAll(wordPressDirectory, os.FileMode(defaultDirPermissions))
	if err != nil {
		return "", err
	}

	return wordPressDirectory, err
}

// getInstalledWordPressPlugins Returns list of installed plugins and whether default plugins are still present.
func (s *Site) getInstalledWordPressPlugins(consoleOutput *console.Console) (pluginList []string, hasDefaultPlugins bool, err error) {
	commands := []string{
		"plugin",
		"list",
		"--format=json",
	}

	hasDefaultPlugins = false

	_, commandOutput, err := s.WPCli(commands, false, consoleOutput)
	if err != nil {
		return []string{}, true, err
	}

	rawPlugins := []PluginInfo{}
	plugins := []string{}

	err = json.Unmarshal([]byte(commandOutput), &rawPlugins)
	if err != nil {
		return []string{}, true, err
	}

	for _, plugin := range rawPlugins {
		if plugin.Status != "dropin" &&
			plugin.Status != "must-use" &&
			plugin.Name != s.settings.Get("name") &&
			plugin.Name != "hello" &&
			plugin.Name != "akismet" {
			plugins = append(plugins, plugin.Name)
		}

		if plugin.Name == "hello" ||
			plugin.Name == "akismet" {
			hasDefaultPlugins = true
		}
	}

	return plugins, hasDefaultPlugins, nil
}

func (s *Site) getWordPressMounts(appDir string) ([]mount.Mount, error) {
	appVolumes := []mount.Mount{
		{ // The root directory of the WordPress site
			Type:   mount.TypeBind,
			Source: appDir,
			Target: "/var/www/html",
		},
		{ // Kana's primary site directory (used for temp files such as DB import and export)
			Type:   mount.TypeBind,
			Source: s.settings.Get("siteDirectory"),
			Target: "/Site",
		},
	}

	wpContentDir := "/var/www/html/wp-content"

	if s.settings.Get("type") == "plugin" {
		err := os.MkdirAll(
			filepath.Join(
				appDir,
				"wp-content",
				"plugins",
				s.settings.Get("name")),
			os.FileMode(defaultDirPermissions))
		if err != nil {
			return appVolumes, err
		}

		appVolumes = append(appVolumes, mount.Mount{ // Map's the user's working directory as a plugin
			Type:   mount.TypeBind,
			Source: s.settings.Get("workingDirectory"),
			Target: filepath.Join(wpContentDir, "plugins", s.settings.Get("name")),
		})
	}

	if s.settings.Get("type") == "theme" {
		err := os.MkdirAll(
			filepath.Join(appDir,
				"wp-content",
				"themes",
				s.settings.Get("name")),
			os.FileMode(defaultDirPermissions))
		if err != nil {
			return appVolumes, err
		}

		appVolumes = append(appVolumes, mount.Mount{ // Map's the user's working directory as a theme
			Type:   mount.TypeBind,
			Source: s.settings.Get("workingDirectory"),
			Target: filepath.Join(wpContentDir, "themes", s.settings.Get("name")),
		})
	}

	return appVolumes, nil
}

func (s *Site) getWordPressContainer(appVolumes []mount.Mount, appContainers []docker.ContainerConfig) []docker.ContainerConfig {
	hostRule := fmt.Sprintf("Host(`%[1]s`)", s.settings.GetDomain())

	envVars := []string{
		"IS_KANA_ENVIRONMENT=true",
	}

	isUsingSQLite, err := s.isUsingSQLite()
	if err != nil {
		return appContainers
	}

	if isUsingSQLite {
		envVars = append(envVars, "KANA_SQLITE=true")
	} else {
		envVars = append(envVars,
			fmt.Sprintf("WORDPRESS_DB_HOST=kana-%s-database", s.settings.Get("name")),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
			"WORDPRESS_ADMIN_USER=admin")
	}

	wordPressContainer := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana-%s-wordpress", s.settings.Get("name")),
		Image:       fmt.Sprintf("wordpress:php%s", s.settings.Get("php")),
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana-%s-wordpress", s.settings.Get("name")),
		Env:         envVars,
		Labels: map[string]string{
			"traefik.enable": "true",
			"kana.type":      "wordpress",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", s.settings.Get("name")): "web",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", s.settings.Get("name")):        hostRule,
			fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", s.settings.Get("name")):      "websecure",
			fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", s.settings.Get("name")):             hostRule,
			fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", s.settings.Get("name")):              "true",
			"kana.site": s.settings.Get("name"),
		},
		Volumes: appVolumes,
	}

	if s.settings.GetBool("AutomaticLogin") {
		wordPressContainer.Env = append(wordPressContainer.Env, "KANA_ADMIN_LOGIN=true")
	}

	if s.settings.GetBool("WPDebug") {
		wordPressContainer.Env = append(wordPressContainer.Env, "WORDPRESS_DEBUG=1")
	}

	extraConfig := fmt.Sprintf("WORDPRESS_CONFIG_EXTRA=define( 'WP_ENVIRONMENT_TYPE', '%s' );", s.settings.Get("environment"))

	if s.settings.GetBool("ScriptDebug") {
		extraConfig += "define( 'SCRIPT_DEBUG', true );"
	}

	wordPressContainer.Env = append(wordPressContainer.Env, extraConfig)

	appContainers = append(appContainers, wordPressContainer)

	return appContainers
}

// getWordPressContainers returns an array of strings containing the container names for the site.
func (s *Site) getWordPressContainers() []string {
	return []string{
		fmt.Sprintf("kana-%s-database", s.settings.Get("name")),
		fmt.Sprintf("kana-%s-wordpress", s.settings.Get("name")),
		fmt.Sprintf("kana-%s-phpmyadmin", s.settings.Get("name")),
		fmt.Sprintf("kana-%s-mailpit", s.settings.Get("name")),
	}
}

func (s *Site) activateProject(consoleOutput *console.Console) error {
	if s.settings.GetBool("Activate") && s.settings.Get("type") != "site" {
		consoleOutput.Println(
			fmt.Sprintf("Activating %s:  %s",
				s.settings.Get("type"),
				consoleOutput.Bold(consoleOutput.Blue(s.settings.Get("name")))))

		setupCommand := []string{
			s.settings.Get("type"),
			"activate",
			s.settings.Get("name"),
		}

		code, _, err := s.WPCli(setupCommand, false, consoleOutput)
		if err != nil {
			return err
		}

		if code != 0 {
			consoleOutput.Warn(
				fmt.Sprintf(
					"Unable to activate %s: %s.",
					s.settings.Get("type"),
					consoleOutput.Bold(consoleOutput.Blue(s.settings.Get("name")))))
		}
	}

	return nil
}

func (s *Site) activateTheme(consoleOutput *console.Console) error {
	if s.settings.Get("type") == "theme" || s.settings.Get("theme") == "" {
		return nil
	}

	consoleOutput.Println(fmt.Sprintf("Installing default theme:  %s", consoleOutput.Bold(consoleOutput.Blue(s.settings.Get("theme")))))

	setupCommand := []string{
		"theme",
		"install",
		"--activate",
		s.settings.Get("theme"),
	}

	code, _, err := s.WPCli(setupCommand, false, consoleOutput)
	if err != nil {
		return err
	}

	if code != 0 {
		consoleOutput.Warn(fmt.Sprintf("Unable to install theme: %s.", consoleOutput.Bold(consoleOutput.Blue(s.settings.Get("theme")))))
	}

	return nil
}

// installDefaultPlugins Installs a list of WordPress plugins.
func (s *Site) installDefaultPlugins(consoleOutput *console.Console) error {
	installedPlugins, _, err := s.getInstalledWordPressPlugins(consoleOutput)
	if err != nil {
		return err
	}

	for _, plugin := range s.settings.GetSlice("plugins") {
		setupCommand := []string{
			"plugin",
			"install",
			"--activate",
			plugin,
		}

		// Don't  try to reinstall the plugin if it is already installed
		for _, installedPlugin := range installedPlugins {
			if installedPlugin == plugin {
				setupCommand = []string{}
			}
		}

		if len(setupCommand) > 0 {
			consoleOutput.Println(fmt.Sprintf("Installing plugin:  %s", consoleOutput.Bold(consoleOutput.Blue(plugin))))

			code, _, err := s.WPCli(setupCommand, false, consoleOutput)
			if err != nil {
				return err
			}

			if code != 0 {
				consoleOutput.Warn(fmt.Sprintf("Unable to install plugin: %s.", consoleOutput.Bold(consoleOutput.Blue(plugin))))
			}
		}
	}

	return nil
}

// installKanaPlugin installs the Kana development plugin.
func (s *Site) installKanaPlugin() error {
	wordPressDirectory, err := s.getWordPressDirectory()
	if err != nil {
		return err
	}

	return settings.EnsureKanaPlugin(wordPressDirectory, s.settings.Get("version"), s.settings.Get("name"))
}

// installWordPress Installs and configures WordPress core.
func (s *Site) installWordPress(consoleOutput *console.Console) error {
	checkCommand := []string{
		"option",
		"get",
		"siteurl",
	}

	code, checkURL, err := s.WPCli(checkCommand, false, consoleOutput)

	if err != nil || code != 0 {
		consoleOutput.Println("Finishing WordPress setup.")

		installCommand := "install"

		if s.settings.Get("multisite") != "none" {
			installCommand = "multisite-install"
		}

		setupCommand := []string{
			"core",
			installCommand,
			fmt.Sprintf("--url=%s", s.settings.GetURL()),
			fmt.Sprintf("--title=Kana Development %s: %s", s.settings.Get("type"), s.settings.Get("name")),
			fmt.Sprintf("--admin_user=%s", s.settings.Get("adminUser")),
			fmt.Sprintf("--admin_password=%s", s.settings.Get("adminPassword")),
			fmt.Sprintf("--admin_email=%s", s.settings.Get("adminEmail")),
		}

		if installCommand == "multisite-install" {
			if s.settings.Get("multisite") == "subdomain" {
				setupCommand = append(setupCommand, "--subdomains")
			}

			err = s.writeHtaccess()
			if err != nil {
				return err
			}
		}

		var output string

		code, output, err = s.WPCli(setupCommand, false, consoleOutput)
		if err != nil || code != 0 {
			return fmt.Errorf("installation of WordPress failed: %s", output)
		}
	} else if strings.TrimSpace(checkURL) != s.settings.GetURL() {
		consoleOutput.Println("The SSL config has changed. Updating the site URL accordingly.")

		// update the home and siteurl to ensure correct ssl usage
		options := []string{
			"siteurl",
			"home",
		}

		for _, option := range options {
			setSiteURLCommand := []string{
				"option",
				"update",
				option,
				s.settings.GetURL(),
			}

			code, _, err = s.WPCli(setSiteURLCommand, false, consoleOutput)
			if err != nil || code != 0 {
				return fmt.Errorf("installation of WordPress failed: %s", err.Error())
			}
		}
	}

	return nil
}

// startWordPress Starts the WordPress containers.
func (s *Site) startWordPress(consoleOutput *console.Console) error {
	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	appDir, databaseDir, err := s.getDirectories()
	if err != nil {
		return err
	}

	// Replace wp-config.php with the container's file
	_, err = os.Stat(filepath.Join(appDir, "wp-config.php"))
	if err == nil {
		os.Remove(filepath.Join(appDir, "wp-config.php"))
	}

	appVolumes, err := s.getWordPressMounts(appDir)
	if err != nil {
		return err
	}

	var appContainers []docker.ContainerConfig

	appContainers = s.getDatabaseContainer(databaseDir, appContainers)
	appContainers = s.getWordPressContainer(appVolumes, appContainers)

	for i := range appContainers {
		err := s.startContainer(&appContainers[i], true, true, consoleOutput)
		if err != nil {
			return err
		}
	}

	return s.verifyDatabase(consoleOutput) // verify the database is ready for connections. On slow filesystems this can take a few seconds.
}

// resetWPFilePermissions Ensures the www-data user owns the WordPress directory.
func (s *Site) resetWPFilePermissions() error {
	if runtime.GOOS == "linux" {
		return nil
	}

	_, err := s.dockerClient.ContainerExec(
		fmt.Sprintf("kana-%s-wordpress", s.settings.Get("name")),
		true,
		[]string{"chown -R www-data:www-data /var/www/html"})

	return err
}

// stopWordPress Stops the site in docker, destroying the containers when they close.
func (s *Site) stopWordPress() error {
	wordPressContainers := s.getWordPressContainers()

	for _, wordPressContainer := range wordPressContainers {
		_, err := s.dockerClient.ContainerStop(wordPressContainer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Site) writeHtaccess() error {
	wordPressDirectory, err := s.getWordPressDirectory()
	if err != nil {
		return err
	}

	_, filePerms := settings.GetDefaultFilePermissions()
	htaccessContents := settings.GetHtaccess(s.settings.Get("multisite"))

	return os.WriteFile(filepath.Join(wordPressDirectory, ".htaccess"), []byte(htaccessContents), os.FileMode(filePerms))
}
