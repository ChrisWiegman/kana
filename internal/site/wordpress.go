package site

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/ChrisWiegman/kana-cli/pkg/docker"

	"github.com/docker/docker/api/types/mount"
	"github.com/logrusorgru/aurora/v4"
)

type PluginInfo struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Update  string `json:"update"`
	Version string `json:"version"`
}

var defaultDirPermissions = 0750

// RunWPCli Runs a wp-cli command returning it's output and any errors
func (s *Site) RunWPCli(command []string) (statusCode int64, output string, err error) {
	appDir := path.Join(s.Settings.SiteDirectory, "app")

	if s.isLocalSite() {
		appDir, err = s.getLocalAppDir()
		if err != nil {
			return 1, "", err
		}
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
		Name:        fmt.Sprintf("kana_%s_wordpress_cli", s.Settings.Name),
		Image:       fmt.Sprintf("wordpress:cli-php%s", s.Settings.PHP),
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana_%s_wordpress_cli", s.Settings.Name),
		Command:     fullCommand,
		Env: []string{
			fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.Settings.Name),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
		},
		Labels: map[string]string{
			"kana.site": s.Settings.Name,
		},
		Volumes: appVolumes,
	}

	err = s.dockerClient.EnsureImage(container.Image)
	if err != nil {
		return 1, "", err
	}

	code, output, err := s.dockerClient.ContainerRunAndClean(&container)
	if err != nil {
		return code, "", err
	}

	return code, output, nil
}

// getInstalledWordPressPlugins Returns a list of the plugins that have been installed on the site
func (s *Site) getInstalledWordPressPlugins() ([]string, error) {
	commands := []string{
		"plugin",
		"list",
		"--format=json",
	}

	_, commandOutput, err := s.RunWPCli(commands)
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
		if plugin.Status != "dropin" && plugin.Name != s.Settings.Name && plugin.Name != "hello" && plugin.Name != "akismet" {
			plugins = append(plugins, plugin.Name)
		}
	}

	return plugins, nil
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

// getWordPressContainers returns an array of strings containing the container names for the site
func (s *Site) getWordPressContainers() []string {
	return []string{
		fmt.Sprintf("kana_%s_database", s.Settings.Name),
		fmt.Sprintf("kana_%s_wordpress", s.Settings.Name),
		fmt.Sprintf("kana_%s_phpmyadmin", s.Settings.Name),
	}
}

// installDefaultPlugins Installs a list of WordPress plugins
func (s *Site) installDefaultPlugins() error {
	installedPlugins, err := s.getInstalledWordPressPlugins()
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
			console.Println(fmt.Sprintf("Installing plugin:  %s", aurora.Bold(aurora.Blue(plugin))))

			setupCommand := []string{
				"plugin",
				"install",
				"--activate",
				plugin,
			}

			code, _, err := s.RunWPCli(setupCommand)
			if err != nil {
				return err
			}

			if code != 0 {
				console.Warn(fmt.Sprintf("Unable to install plugin: %s.", aurora.Bold(aurora.Blue(plugin))))
			}
		}
	}

	return nil
}

// installWordPress Installs and configures WordPress core
func (s *Site) installWordPress() error {
	checkCommand := []string{
		"core",
		"is-installed",
	}

	code, _, err := s.RunWPCli(checkCommand)

	if err != nil || code != 0 {
		console.Println("Finishing WordPress setup.")

		setupCommand := []string{
			"core",
			"install",
			fmt.Sprintf("--url=%s", s.getSiteURL(false)),
			fmt.Sprintf("--title=Kana Development %s: %s", s.Settings.Type, s.Settings.Name),
			fmt.Sprintf("--admin_user=%s", s.Settings.AdminUsername),
			fmt.Sprintf("--admin_password=%s", s.Settings.AdminPassword),
			fmt.Sprintf("--admin_email=%s", s.Settings.AdminEmail),
		}

		code, _, err = s.RunWPCli(setupCommand)
		if err != nil || code != 0 {
			return fmt.Errorf("installation of WordPress failed: %s", err.Error())
		}
	}

	return nil
}

// startWordPress Starts the WordPress containers
func (s *Site) startWordPress() error {
	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	appDir := path.Join(s.Settings.SiteDirectory, "app")
	databaseDir := path.Join(s.Settings.SiteDirectory, "database")

	if s.isLocalSite() {
		appDir, err = s.getLocalAppDir()
		if err != nil {
			return err
		}

		// Replace wp-config.php with the container's file
		_, err = os.Stat(path.Join(appDir, "wp-config.php"))
		if err == nil {
			os.Remove(path.Join(appDir, "wp-config.php"))
		}
	}

	err = os.MkdirAll(appDir, os.FileMode(defaultDirPermissions))
	if err != nil {
		return err
	}

	err = os.MkdirAll(databaseDir, os.FileMode(defaultDirPermissions))
	if err != nil {
		return err
	}

	appVolumes, err := s.getMounts(appDir)
	if err != nil {
		return err
	}

	wordPressContainers := []docker.ContainerConfig{
		{
			Name:        fmt.Sprintf("kana_%s_database", s.Settings.Name),
			Image:       "mariadb",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_database", s.Settings.Name),
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
		},
		{
			Name:        fmt.Sprintf("kana_%s_wordpress", s.Settings.Name),
			Image:       fmt.Sprintf("wordpress:php%s", s.Settings.PHP),
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_wordpress", s.Settings.Name),
			Env: []string{
				fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.Settings.Name),
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
		},
	}

	wordPressContainers = s.getPhpMyAdminContainer(databaseDir, wordPressContainers)

	for i := range wordPressContainers {
		err := s.dockerClient.EnsureImage(wordPressContainers[i].Image)
		if err != nil {
			return err
		}
		_, err = s.dockerClient.ContainerRun(&wordPressContainers[i], true, true)
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

func (s *Site) getPhpMyAdminContainer(databaseDir string, wordPressContainers []docker.ContainerConfig) []docker.ContainerConfig {
	if s.Settings.PhpMyAdmin {
		phpMyAdminContainer := docker.ContainerConfig{
			Name:        fmt.Sprintf("kana_%s_phpmyadmin", s.Settings.Name),
			Image:       "phpmyadmin",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_phpmyadmin", s.Settings.Name),
			Env: []string{
				"MYSQL_ROOT_PASSWORD=password",
				fmt.Sprintf("PMA_HOST=kana_%s_database", s.Settings.Name),
				"PMA_USER=wordpress",
				"PMA_PASSWORD=wordpress",
			},
			Volumes: []mount.Mount{
				{ // Maps a database folder to the MySQL container for persistence
					Type:   mount.TypeBind,
					Source: databaseDir,
					Target: "/var/lib/mysql",
				},
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

		wordPressContainers = append(wordPressContainers, phpMyAdminContainer)
	}

	return wordPressContainers
}
