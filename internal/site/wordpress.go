package site

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/ChrisWiegman/kana-cli/pkg/docker"

	"github.com/docker/docker/api/types/mount"
)

type PluginInfo struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Update  string `json:"update"`
	Version string `json:"version"`
}

var defaultDirPermissions = 0750

// RunWPCli Runs a wp-cli command returning it's output and any errors
func (s *Site) RunWPCli(command []string, consoleOutput *console.Console) (statusCode int64, output string, err error) {
	appDir, err := s.getAppDirectory(consoleOutput)
	if err != nil {
		return 1, "", err
	}

	appVolumes, err := s.getMounts(appDir)
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

	err = s.dockerClient.EnsureImage(container.Image, consoleOutput)
	if err != nil {
		return 1, "", err
	}

	code, output, err := s.dockerClient.ContainerRunAndClean(&container)
	if err != nil {
		return code, "", err
	}

	return code, output, nil
}

func (s *Site) getAppDirectory(consoleOutput *console.Console) (string, error) {
	var err error
	appDir := path.Join(s.Settings.SiteDirectory, "app")

	if s.isLocalSite(consoleOutput) {
		appDir, err = s.getLocalAppDir()
		if err != nil {
			return "", err
		}
	}

	return appDir, err
}

func (s *Site) getDatabaseContainer(databaseDir string, appContainers []docker.ContainerConfig) []docker.ContainerConfig {
	databaseContainer := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana-%s-database", s.Settings.Name),
		Image:       "mariadb:10",
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana-%s-database", s.Settings.Name),
		Ports: []docker.ExposedPorts{
			{Port: "3306", Protocol: "tcp"},
		},
		Env: []string{
			"MARIADB_ROOT_PASSWORD=password",
			"MARIADB_DATABASE=wordpress",
			"MARIADB_USER=wordpress",
			"MARIADB_PASSWORD=wordpress",
		},
		Labels: map[string]string{
			"kana.site": s.Settings.Name,
		},
		Volumes: []mount.Mount{
			{ // Maps a database folder to the MySQL container for persistence
				Type:   mount.TypeBind,
				Source: databaseDir,
				Target: "/var/lib/mysql",
			},
		},
	}

	appContainers = append(appContainers, databaseContainer)

	return appContainers
}

// getDirectories Returns the correct appDir and databaseDir for the current site
func (s *Site) getDirectories(consoleOutput *console.Console) (appDir, databaseDir string, err error) {
	appDir, err = s.getAppDirectory(consoleOutput)
	if err != nil {
		return "", "", err
	}

	databaseDir = path.Join(s.Settings.SiteDirectory, "database")

	err = os.MkdirAll(appDir, os.FileMode(defaultDirPermissions))
	if err != nil {
		return "", "", err
	}

	err = os.MkdirAll(databaseDir, os.FileMode(defaultDirPermissions))
	return appDir, databaseDir, err
}

// getInstalledWordPressPlugins Returns a list of the plugins that have been installed on the site
func (s *Site) getInstalledWordPressPlugins(consoleOutput *console.Console) ([]string, error) {
	commands := []string{
		"plugin",
		"list",
		"--format=json",
	}

	_, commandOutput, err := s.RunWPCli(commands, consoleOutput)
	if err != nil {
		return []string{}, err
	}

	rawPlugins := []PluginInfo{}
	plugins := []string{}

	err = json.Unmarshal([]byte(commandOutput), &rawPlugins)
	if err != nil {
		return []string{}, err
	}

	for _, plugin := range rawPlugins {
		if plugin.Status != "dropin" &&
			plugin.Status != "must-use" &&
			plugin.Name != s.Settings.Name &&
			plugin.Name != "hello" &&
			plugin.Name != "akismet" {
			plugins = append(plugins, plugin.Name)
		}
	}

	return plugins, nil
}

func (s *Site) getMailpitContainer() docker.ContainerConfig {
	mailpitContainer := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana-%s-mailpit", s.Settings.Name),
		Image:       "axllent/mailpit",
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana-%s-mailpit", s.Settings.Name),
		Env:         []string{},
		Volumes:     []mount.Mount{},
		Ports: []docker.ExposedPorts{
			{Port: "8025", Protocol: "tcp"},
			{Port: "1025", Protocol: "tcp"},
		},
		Labels: map[string]string{
			"traefik.enable": "true",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-%s-http.entrypoints", s.Settings.Name, "mailpit"): "web",
			fmt.Sprintf(
				"traefik.http.routers.wordpress-%s-%s-http.rule",
				s.Settings.Name,
				"mailpit"): fmt.Sprintf(
				"Host(`%s-%s`)",
				"mailpit",
				s.Settings.SiteDomain),
			fmt.Sprintf("traefik.http.routers.wordpress-%s-%s.entrypoints", s.Settings.Name, "mailpit"): "websecure",
			fmt.Sprintf(
				"traefik.http.routers.wordpress-%s-%s.rule",
				s.Settings.Name,
				"mailpit"): fmt.Sprintf(
				"Host(`%s-%s`)",
				"mailpit",
				s.Settings.SiteDomain),
			fmt.Sprintf("traefik.http.services.%s-http-svc.loadbalancer.server.port", "mailpit"): "8025",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-%s.tls", s.Settings.Name, "mailpit"):  "true",
			"kana.site": s.Settings.Name,
		},
	}

	return mailpitContainer
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
				s.Settings.WorkingDirectory,
				"wordpress",
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
			path.Join(s.Settings.WorkingDirectory,
				"wordpress",
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

func (s *Site) getPhpMyAdminContainer() docker.ContainerConfig {
	phpMyAdminContainer := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana-%s-phpmyadmin", s.Settings.Name),
		Image:       "phpmyadmin",
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana-%s-phpmyadmin", s.Settings.Name),
		Env: []string{
			"MYSQL_ROOT_PASSWORD=password",
			fmt.Sprintf("PMA_HOST=kana-%s-database", s.Settings.Name),
			"PMA_USER=wordpress",
			"PMA_PASSWORD=wordpress",
		},
		Labels: map[string]string{
			"traefik.enable": "true",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-%s-http.entrypoints", s.Settings.Name, "phpmyadmin"): "web",
			fmt.Sprintf(
				"traefik.http.routers.wordpress-%s-%s-http.rule",
				s.Settings.Name,
				"phpmyadmin"): fmt.Sprintf(
				"Host(`%s-%s`)",
				"phpmyadmin",
				s.Settings.SiteDomain),
			fmt.Sprintf("traefik.http.routers.wordpress-%s-%s.entrypoints", s.Settings.Name, "phpmyadmin"): "websecure",
			fmt.Sprintf(
				"traefik.http.routers.wordpress-%s-%s.rule",
				s.Settings.Name,
				"phpmyadmin"): fmt.Sprintf(
				"Host(`%s-%s`)",
				"phpmyadmin",
				s.Settings.SiteDomain),
			fmt.Sprintf("traefik.http.routers.wordpress-%s-%s.tls", s.Settings.Name, "phpmyadmin"): "true",
			"kana.site": s.Settings.Name,
		},
	}

	return phpMyAdminContainer
}

func (s *Site) getWordPressContainer(appVolumes []mount.Mount, appContainers []docker.ContainerConfig) []docker.ContainerConfig {
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
			fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", s.Settings.Name): "web",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", s.Settings.Name):        fmt.Sprintf("Host(`%s`)", s.Settings.SiteDomain),
			fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", s.Settings.Name):      "websecure",
			fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", s.Settings.Name):             fmt.Sprintf("Host(`%s`)", s.Settings.SiteDomain),
			fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", s.Settings.Name):              "true",
			"kana.site": s.Settings.Name,
		},
		Volumes: appVolumes,
	}

	if s.Settings.WPDebug {
		wordPressContainer.Env = append(wordPressContainer.Env, "WORDPRESS_DEBUG=1")
	}

	appContainers = append(appContainers, wordPressContainer)

	return appContainers
}

// getWordPressContainers returns an array of strings containing the container names for the site
func (s *Site) getWordPressContainers() []string {
	return []string{
		fmt.Sprintf("kana-%s-database", s.Settings.Name),
		fmt.Sprintf("kana-%s-wordpress", s.Settings.Name),
		fmt.Sprintf("kana-%s-phpmyadmin", s.Settings.Name),
		fmt.Sprintf("kana-%s-mailpit", s.Settings.Name),
	}
}

// installDefaultPlugins Installs a list of WordPress plugins
func (s *Site) installDefaultPlugins(consoleOutput *console.Console) error {
	installedPlugins, err := s.getInstalledWordPressPlugins(consoleOutput)
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

			code, _, err := s.RunWPCli(setupCommand, consoleOutput)
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

// installKanaPlugin installs the Kana development plugin
func (s *Site) installKanaPlugin(consoleOutput *console.Console) error {
	appDir, err := s.getAppDirectory(consoleOutput)
	if err != nil {
		return err
	}

	return s.Settings.EnsureKanaPlugin(appDir)
}

// installWordPress Installs and configures WordPress core
func (s *Site) installWordPress(consoleOutput *console.Console) error {
	checkCommand := []string{
		"option",
		"get",
		"siteurl",
	}

	code, checkURL, err := s.RunWPCli(checkCommand, consoleOutput)

	if err != nil || code != 0 {
		consoleOutput.Println("Finishing WordPress setup.")

		setupCommand := []string{
			"core",
			"install",
			fmt.Sprintf("--url=%s", s.Settings.URL),
			fmt.Sprintf("--title=Kana Development %s: %s", s.Settings.Type, s.Settings.Name),
			fmt.Sprintf("--admin_user=%s", s.Settings.AdminUsername),
			fmt.Sprintf("--admin_password=%s", s.Settings.AdminPassword),
			fmt.Sprintf("--admin_email=%s", s.Settings.AdminEmail),
		}

		code, _, err = s.RunWPCli(setupCommand, consoleOutput)
		if err != nil || code != 0 {
			return fmt.Errorf("installation of WordPress failed: %s", err.Error())
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

			code, _, err = s.RunWPCli(setSiteURLCommand, consoleOutput)
			if err != nil || code != 0 {
				return fmt.Errorf("installation of WordPress failed: %s", err.Error())
			}
		}
	}

	return nil
}

// startContainer Starts a given container configuration
func (s *Site) startContainer(container *docker.ContainerConfig, randomPorts, localUser bool, consoleOutput *console.Console) error {
	err := s.dockerClient.EnsureImage(container.Image, consoleOutput)
	if err != nil {
		return err
	}
	_, err = s.dockerClient.ContainerRun(container, randomPorts, localUser)

	return err
}

// startMailpit Starts the Mailpit container
func (s *Site) startMailpit(consoleOutput *console.Console) error {
	mailpitContainer := s.getMailpitContainer()

	return s.startContainer(&mailpitContainer, true, true, consoleOutput)
}

// startPHPMyAdmin Starts the PhpMyAdmin container
func (s *Site) startPHPMyAdmin(consoleOutput *console.Console) error {
	phpMyAdminContainer := s.getPhpMyAdminContainer()

	return s.startContainer(&phpMyAdminContainer, true, false, consoleOutput)
}

// startWordPress Starts the WordPress containers
func (s *Site) startWordPress(consoleOutput *console.Console) error {
	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	appDir, databaseDir, err := s.getDirectories(consoleOutput)
	if err != nil {
		return err
	}

	if s.isLocalSite(consoleOutput) {
		// Replace wp-config.php with the container's file
		_, err = os.Stat(path.Join(appDir, "wp-config.php"))
		if err == nil {
			os.Remove(path.Join(appDir, "wp-config.php"))
		}
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

// stopWordPress Stops the site in docker, destroying the containers when they close
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
