package site

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/docker"

	"github.com/docker/docker/api/types/mount"
)

type CurrentConfig struct {
	Type   string
	Local  bool
	Xdebug bool
}

type PluginInfo struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Update  string `json:"update"`
	Version string `json:"version"`
}

// RunWPCli Runs a wp-cli command returning it's output and any errors
func (s *Site) RunWPCli(command []string) (string, error) {

	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return "", err
	}

	siteDir := path.Join(s.Config.Directories.App, "sites", s.Config.Site.Name)
	appDir := path.Join(siteDir, "app")
	runningConfig := s.getRunningConfig()

	if runningConfig.Local {
		appDir, err = getLocalAppDir()
		if err != nil {
			return "", err
		}
	}

	appVolumes, err := s.getMounts(s.Config.Directories.Site, appDir, runningConfig.Type)
	if err != nil {
		return "", err
	}

	fullCommand := []string{
		"wp",
		"--path=/var/www/html",
	}

	fullCommand = append(fullCommand, command...)

	container := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana_%s_wordpress_cli", s.Config.Site.Name),
		Image:       fmt.Sprintf("wordpress:cli-php%s", s.Config.Site.PHP),
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana_%s_wordpress_cli", s.Config.Site.Name),
		Command:     fullCommand,
		Env: []string{
			fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.Config.Site.Name),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
		},
		Labels: map[string]string{
			"kana.site": s.Config.Site.Name,
		},
		Volumes: appVolumes,
	}

	err = s.dockerClient.EnsureImage(container.Image)
	if err != nil {
		return "", err
	}

	_, output, err := s.dockerClient.ContainerRunAndClean(container)
	if err != nil {
		return "", err
	}

	return output, nil
}

// getInstalledWordPressPlugins Returns a list of the plugins that have been installed on the site
func (s *Site) getInstalledWordPressPlugins() ([]string, error) {

	commands := []string{
		"plugin",
		"list",
		"--format=json",
	}

	commandOutput, err := s.RunWPCli(commands)
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

		if plugin.Status != "dropin" && plugin.Name != s.Config.Site.Name && plugin.Name != "hello" && plugin.Name != "akismet" {
			plugins = append(plugins, plugin.Name)
		}
	}

	return plugins, nil
}

func (s *Site) getMounts(siteDir, appDir, siteType string) ([]mount.Mount, error) {

	appVolumes := []mount.Mount{
		{ // The root directory of the WordPress site
			Type:   mount.TypeBind,
			Source: appDir,
			Target: "/var/www/html",
		},
		{ // Kana's primary site directory (used for temp files such as DB import and export)
			Type:   mount.TypeBind,
			Source: siteDir,
			Target: "/Site",
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		return appVolumes, err
	}

	if siteType == "plugin" {
		appVolumes = append(appVolumes, mount.Mount{ // Map's the user's working directory as a plugin
			Type:   mount.TypeBind,
			Source: cwd,
			Target: path.Join("/var/www/html", "wp-content", "plugins", s.Config.Site.Name),
		})
	}

	if siteType == "theme" {
		appVolumes = append(appVolumes, mount.Mount{ // Map's the user's working directory as a theme
			Type:   mount.TypeBind,
			Source: cwd,
			Target: path.Join("/var/www/html", "wp-content", "themes", s.Config.Site.Name),
		})
	}

	return appVolumes, nil
}

// getWordPressContainers returns an array of strings containing the container names for the site
func (s *Site) getWordPressContainers() []string {

	return []string{
		fmt.Sprintf("kana_%s_database", s.Config.Site.Name),
		fmt.Sprintf("kana_%s_wordpress", s.Config.Site.Name),
	}
}

// installDefaultPlugins Installs a list of WordPress plugins
func (s *Site) installDefaultPlugins() error {

	for _, plugin := range s.Config.Site.Plugins {

		setupCommand := []string{
			"plugin",
			"install",
			"--activate",
			plugin,
		}

		_, err := s.RunWPCli(setupCommand)
		if err != nil {
			return err
		}
	}

	return nil
}

// installWordPress Installs and configures WordPress core
func (s *Site) installWordPress() error {

	console.Println("Finishing WordPress setup...")

	setupCommand := []string{
		"core",
		"install",
		fmt.Sprintf("--url=%s", s.getSiteURL(false)),
		fmt.Sprintf("--title=Kana Development %s: %s", s.Config.Site.Type, s.Config.Site.Name),
		fmt.Sprintf("--admin_user=%s", s.Config.App.AdminUsername),
		fmt.Sprintf("--admin_password=%s", s.Config.App.AdminPassword),
		fmt.Sprintf("--admin_email=%s", s.Config.App.AdminEmail),
	}

	_, err := s.RunWPCli(setupCommand)
	return err
}

// startWordPress Starts the WordPress containers
func (s *Site) startWordPress() error {

	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	appDir := path.Join(s.Config.Directories.Site, "app")
	databaseDir := path.Join(s.Config.Directories.Site, "database")

	if s.isLocalSite() {

		appDir, err = getLocalAppDir()
		if err != nil {
			return err
		}

		// Replace wp-config.php with the container's file
		_, err := os.Stat(path.Join(appDir, "wp-config.php"))
		if err == nil {
			os.Remove(path.Join(appDir, "wp-config.php"))
		}
	}

	if err := os.MkdirAll(appDir, 0750); err != nil {
		return err
	}

	if err := os.MkdirAll(databaseDir, 0750); err != nil {
		return err
	}

	appVolumes, err := s.getMounts(s.Config.Directories.Working, appDir, s.Config.Site.Type)
	if err != nil {
		return err
	}

	wordPressContainers := []docker.ContainerConfig{
		{
			Name:        fmt.Sprintf("kana_%s_database", s.Config.Site.Name),
			Image:       "mariadb",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_database", s.Config.Site.Name),
			Env: []string{
				"MARIADB_ROOT_PASSWORD=password",
				"MARIADB_DATABASE=wordpress",
				"MARIADB_USER=wordpress",
				"MARIADB_PASSWORD=wordpress",
			},
			Labels: map[string]string{
				"kana.site": s.Config.Site.Name,
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
			Name:        fmt.Sprintf("kana_%s_wordpress", s.Config.Site.Name),
			Image:       fmt.Sprintf("wordpress:php%s", s.Config.Site.PHP),
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_wordpress", s.Config.Site.Name),
			Env: []string{
				fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.Config.Site.Name),
				"WORDPRESS_DB_USER=wordpress",
				"WORDPRESS_DB_PASSWORD=wordpress",
				"WORDPRESS_DB_NAME=wordpress",
			},
			Labels: map[string]string{
				"traefik.enable": "true",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", s.Config.Site.Name): "web",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", s.Config.Site.Name):        fmt.Sprintf("Host(`%s`)", s.Config.Site.Domain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", s.Config.Site.Name):      "websecure",
				fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", s.Config.Site.Name):             fmt.Sprintf("Host(`%s`)", s.Config.Site.Domain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", s.Config.Site.Name):              "true",
				"kana.site": s.Config.Site.Name,
			},
			Volumes: appVolumes,
		},
	}

	for _, container := range wordPressContainers {

		err := s.dockerClient.EnsureImage(container.Image)
		if err != nil {
			return err
		}

		_, err = s.dockerClient.ContainerRun(container)
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
