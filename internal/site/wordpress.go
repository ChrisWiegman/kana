package site

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/docker"
	"github.com/ChrisWiegman/kana-cli/internal/settings"

	"github.com/docker/docker/api/types/mount"
)

type PluginInfo struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Update  string `json:"update"`
	Version string `json:"version"`
}

var defaultDirPermissions = 0750

// RunWPCli Runs a wp-cli command returning it's output and any errors.
func (s *Site) RunWPCli(command []string, interactive bool, consoleOutput *console.Console) (statusCode int64, output string, err error) {
	mounts := s.dockerClient.ContainerGetMounts(fmt.Sprintf("kana-%s-wordpress", s.Settings.Name))

	for _, mount := range mounts {
		if strings.Contains(mount.Destination, "/var/www/html/wp-content/plugins/") {
			s.Settings.Type = "plugin"
		}

		if strings.Contains(mount.Destination, "/var/www/html/wp-content/themes/") {
			s.Settings.Type = "theme"
		}
	}

	wordPressDirectory, err := s.getWordPressDirectory()
	if err != nil {
		return 1, "", err
	}

	appVolumes, err := s.getMounts(wordPressDirectory)
	if err != nil {
		return 1, "", err
	}

	fullCommand := []string{
		"wp",
		"--path=/var/www/html",
	}

	fullCommand = append(fullCommand, command...)

	container := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana-%s-wordpress_cli", s.Settings.Name),
		Image:       fmt.Sprintf("wordpress:cli-php%s", s.Settings.PHP),
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana-%s-wordpress_cli", s.Settings.Name),
		Command:     fullCommand,
		Env: []string{
			fmt.Sprintf("WORDPRESS_DB_HOST=kana-%s-database", s.Settings.Name),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
			"WORDPRESS_DEBUG=0",
		},
		Labels: map[string]string{
			"kana.site": s.Settings.Name,
		},
		Volumes: appVolumes,
	}

	err = s.dockerClient.EnsureImage(container.Image, s.Settings.ImageUpdateDays, consoleOutput)
	if err != nil {
		return 1, "", err
	}

	code, output, err := s.dockerClient.ContainerRunAndClean(&container, interactive)
	if err != nil {
		return code, "", err
	}

	return code, output, nil
}

func (s *Site) getDatabaseDirectory() (databaseDirectory string, err error) {
	databaseDirectory = path.Join(s.Settings.SiteDirectory, "database")

	err = os.MkdirAll(databaseDirectory, os.FileMode(defaultDirPermissions))
	if err != nil {
		return "", err
	}

	return databaseDirectory, err
}

func (s *Site) getWordPressDirectory() (wordPressDirectory string, err error) {
	wordPressDirectory = path.Join(s.Settings.SiteDirectory, "wordpress")

	if !s.Settings.IsNamedSite {
		wordPressDirectory = s.Settings.WorkingDirectory

		if s.Settings.Type != DefaultType {
			wordPressDirectory = path.Join(s.Settings.WorkingDirectory, "wordpress")
		}
	}

	err = os.MkdirAll(wordPressDirectory, os.FileMode(defaultDirPermissions))
	if err != nil {
		return "", err
	}

	return wordPressDirectory, err
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

// getInstalledWordPressPlugins Returns list of installed plugins and whether default plugins are still present.
func (s *Site) getInstalledWordPressPlugins(consoleOutput *console.Console) (pluginList []string, hasDefaultPlugins bool, err error) {
	commands := []string{
		"plugin",
		"list",
		"--format=json",
	}

	hasDefaultPlugins = false

	_, commandOutput, err := s.RunWPCli(commands, false, consoleOutput)
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
			plugin.Name != s.Settings.Name &&
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

func (s *Site) getMounts(appDir string) ([]mount.Mount, error) {
	appVolumes := []mount.Mount{
		{ // The root directory of the WordPress site
			Type:   mount.TypeBind,
			Source: appDir,
			Target: "/var/www/html",
		},
		{ // Kana's primary site directory (used for temp files such as DB import and export)
			Type:   mount.TypeBind,
			Source: s.Settings.SiteDirectory,
			Target: "/Site",
		},
	}

	if s.Settings.Type == "plugin" {
		err := os.MkdirAll(
			path.Join(
				appDir,
				"wp-content",
				"plugins",
				s.Settings.Name),
			os.FileMode(defaultDirPermissions))
		if err != nil {
			return appVolumes, err
		}

		appVolumes = append(appVolumes, mount.Mount{ // Map's the user's working directory as a plugin
			Type:   mount.TypeBind,
			Source: s.Settings.WorkingDirectory,
			Target: path.Join("/var/www/html", "wp-content", "plugins", s.Settings.Name),
		})
	}

	if s.Settings.Type == "theme" {
		err := os.MkdirAll(
			path.Join(appDir,
				"wp-content",
				"themes",
				s.Settings.Name),
			os.FileMode(defaultDirPermissions))
		if err != nil {
			return appVolumes, err
		}

		appVolumes = append(appVolumes, mount.Mount{ // Map's the user's working directory as a theme
			Type:   mount.TypeBind,
			Source: s.Settings.WorkingDirectory,
			Target: path.Join("/var/www/html", "wp-content", "themes", s.Settings.Name),
		})
	}

	return appVolumes, nil
}

func (s *Site) getWordPressContainer(appVolumes []mount.Mount, appContainers []docker.ContainerConfig) []docker.ContainerConfig {
	hostRule := fmt.Sprintf("HostRegexp(`%[1]s`, `{subdomain:.*}.%[1]s`)", s.Settings.SiteDomain)

	wordPressContainer := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana-%s-wordpress", s.Settings.Name),
		Image:       fmt.Sprintf("wordpress:php%s", s.Settings.PHP),
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana-%s-wordpress", s.Settings.Name),
		Env: []string{
			fmt.Sprintf("WORDPRESS_DB_HOST=kana-%s-database", s.Settings.Name),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
		},
		Labels: map[string]string{
			"traefik.enable": "true",
			"kana.type":      "wordpress",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", s.Settings.Name): "web",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", s.Settings.Name):        hostRule,
			fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", s.Settings.Name):      "websecure",
			fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", s.Settings.Name):             hostRule,
			fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", s.Settings.Name):              "true",
			"kana.site": s.Settings.Name,
		},
		Volumes: appVolumes,
	}

	if s.Settings.WPDebug {
		wordPressContainer.Env = append(wordPressContainer.Env, "WORDPRESS_DEBUG=1")
	}

	extraConfig := fmt.Sprintf("WORDPRESS_CONFIG_EXTRA=define( 'WP_ENVIRONMENT_TYPE', '%s' );", s.Settings.Environment)

	if s.Settings.ScriptDebug {
		extraConfig += "define( 'SCRIPT_DEBUG', true );"
	}

	wordPressContainer.Env = append(wordPressContainer.Env, extraConfig)

	appContainers = append(appContainers, wordPressContainer)

	return appContainers
}

// getWordPressContainers returns an array of strings containing the container names for the site.
func (s *Site) getWordPressContainers() []string {
	return []string{
		fmt.Sprintf("kana-%s-database", s.Settings.Name),
		fmt.Sprintf("kana-%s-wordpress", s.Settings.Name),
		fmt.Sprintf("kana-%s-phpmyadmin", s.Settings.Name),
		fmt.Sprintf("kana-%s-mailpit", s.Settings.Name),
	}
}

func (s *Site) activateProject(consoleOutput *console.Console) error {
	if s.Settings.Activate && s.Settings.Type != "site" {
		consoleOutput.Println(fmt.Sprintf("Activating %s:  %s", s.Settings.Type, consoleOutput.Bold(consoleOutput.Blue(s.Settings.Name))))

		setupCommand := []string{
			s.Settings.Type,
			"activate",
			s.Settings.Name,
		}

		code, _, err := s.RunWPCli(setupCommand, false, consoleOutput)
		if err != nil {
			return err
		}

		if code != 0 {
			consoleOutput.Warn(fmt.Sprintf("Unable to activate %s: %s.", s.Settings.Type, consoleOutput.Bold(consoleOutput.Blue(s.Settings.Name))))
		}
	}

	return nil
}

// installDefaultPlugins Installs a list of WordPress plugins.
func (s *Site) installDefaultPlugins(consoleOutput *console.Console) error {
	installedPlugins, _, err := s.getInstalledWordPressPlugins(consoleOutput)
	if err != nil {
		return err
	}

	for _, plugin := range s.Settings.Plugins {
		installPlugin := true

		for _, installedPlugin := range installedPlugins {
			if installedPlugin == plugin {
				installPlugin = false
			}
		}

		if installPlugin {
			consoleOutput.Println(fmt.Sprintf("Installing plugin:  %s", consoleOutput.Bold(consoleOutput.Blue(plugin))))

			setupCommand := []string{
				"plugin",
				"install",
				"--activate",
				plugin,
			}

			code, _, err := s.RunWPCli(setupCommand, false, consoleOutput)
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

	return s.Settings.EnsureKanaPlugin(wordPressDirectory, s.Settings.Name)
}

// installWordPress Installs and configures WordPress core.
func (s *Site) installWordPress(consoleOutput *console.Console) error {
	checkCommand := []string{
		"option",
		"get",
		"siteurl",
	}

	code, checkURL, err := s.RunWPCli(checkCommand, false, consoleOutput)

	if err != nil || code != 0 {
		consoleOutput.Println("Finishing WordPress setup.")

		installCommand := "install"

		if s.Settings.Multisite != "none" {
			installCommand = "multisite-install"
		}

		setupCommand := []string{
			"core",
			installCommand,
			fmt.Sprintf("--url=%s", s.Settings.URL),
			fmt.Sprintf("--title=Kana Development %s: %s", s.Settings.Type, s.Settings.Name),
			fmt.Sprintf("--admin_user=%s", s.Settings.AdminUsername),
			fmt.Sprintf("--admin_password=%s", s.Settings.AdminPassword),
			fmt.Sprintf("--admin_email=%s", s.Settings.AdminEmail),
		}

		if installCommand == "multisite-install" {
			if s.Settings.Multisite == "subdomain" {
				setupCommand = append(setupCommand, "--subdomains")
			}

			err = s.writeHtaccess()
			if err != nil {
				return err
			}
		}

		var output string

		code, output, err = s.RunWPCli(setupCommand, false, consoleOutput)
		if err != nil || code != 0 {
			return fmt.Errorf("installation of WordPress failed: %s", output)
		}
	} else if strings.TrimSpace(checkURL) != s.Settings.URL {
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
				s.Settings.URL,
			}

			code, _, err = s.RunWPCli(setSiteURLCommand, false, consoleOutput)
			if err != nil || code != 0 {
				return fmt.Errorf("installation of WordPress failed: %s", err.Error())
			}
		}
	}

	return nil
}

// startContainer Starts a given container configuration.
func (s *Site) startContainer(container *docker.ContainerConfig, randomPorts, localUser bool, consoleOutput *console.Console) error {
	err := s.dockerClient.EnsureImage(container.Image, s.Settings.ImageUpdateDays, consoleOutput)
	if err != nil {
		err = s.handleImageError(container, err)
		if err != nil {
			return err
		}
	}
	_, err = s.dockerClient.ContainerRun(container, randomPorts, localUser)

	return err
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
	_, err = os.Stat(path.Join(appDir, "wp-config.php"))
	if err == nil {
		os.Remove(path.Join(appDir, "wp-config.php"))
	}

	appVolumes, err := s.getMounts(appDir)
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

	return nil
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

	_, filePerms := settings.GetDefaultPermissions()
	htaccessContents := s.Settings.GetHtaccess()

	return os.WriteFile(path.Join(wordPressDirectory, ".htaccess"), []byte(htaccessContents), os.FileMode(filePerms))
}
