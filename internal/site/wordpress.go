package site

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/traefik"
	"github.com/docker/docker/api/types/mount"
)

func (s *Site) StopWordPress() error {

	wordPressContainers := []string{
		fmt.Sprintf("kana_%s_database", s.appConfig.SiteName),
		fmt.Sprintf("kana_%s_wordpress", s.appConfig.SiteName),
	}

	for _, wordPressContainer := range wordPressContainers {
		_, err := s.dockerClient.ContainerStop(wordPressContainer)
		if err != nil {
			return err
		}
	}

	traefikClient, err := traefik.NewTraefik(s.appConfig)
	if err != nil {
		return err
	}

	return traefikClient.MaybeStopTraefik()

}

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

func (s *Site) StartWordPress(local, isPlugin, isTheme bool) error {

	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	appDir := path.Join(s.appConfig.SiteDirectory, "app")
	databaseDir := path.Join(s.appConfig.SiteDirectory, "database")

	if local {
		appDir, err = getLocalAppDir()
		if err != nil {
			return err
		}
	}

	if err := os.MkdirAll(appDir, 0750); err != nil {
		return err
	}

	if err := os.MkdirAll(databaseDir, 0750); err != nil {
		return err
	}

	appVolumes := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: appDir,
			Target: "/var/www/html",
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if isPlugin {
		appVolumes = append(appVolumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: cwd,
			Target: path.Join("/var/www/html", "wp-content", "plugins", s.appConfig.SiteName),
		})
	}

	if isTheme {
		appVolumes = append(appVolumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: cwd,
			Target: path.Join("/var/www/html", "wp-content", "themes", s.appConfig.SiteName),
		})
	}

	wordPressContainers := []docker.ContainerConfig{
		{
			Name:        fmt.Sprintf("kana_%s_database", s.appConfig.SiteName),
			Image:       "mariadb",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_database", s.appConfig.SiteName),
			Env: []string{
				"MARIADB_ROOT_PASSWORD=password",
				"MARIADB_DATABASE=wordpress",
				"MARIADB_USER=wordpress",
				"MARIADB_PASSWORD=wordpress",
			},
			Labels: map[string]string{
				"kana.site": s.appConfig.SiteName,
			},
			Volumes: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: databaseDir,
					Target: "/var/lib/mysql",
				},
			},
		},
		{
			Name:        fmt.Sprintf("kana_%s_wordpress", s.appConfig.SiteName),
			Image:       "wordpress",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_wordpress", s.appConfig.SiteName),
			Env: []string{
				fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.appConfig.SiteName),
				"WORDPRESS_DB_USER=wordpress",
				"WORDPRESS_DB_PASSWORD=wordpress",
				"WORDPRESS_DB_NAME=wordpress",
			},
			Labels: map[string]string{
				"traefik.enable": "true",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", s.appConfig.SiteName): "web",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", s.appConfig.SiteName):        fmt.Sprintf("Host(`%s.%s`)", s.appConfig.SiteName, s.appConfig.AppDomain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", s.appConfig.SiteName):      "websecure",
				fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", s.appConfig.SiteName):             fmt.Sprintf("Host(`%s.%s`)", s.appConfig.SiteName, s.appConfig.AppDomain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", s.appConfig.SiteName):              "true",
				"kana.site": s.appConfig.SiteName,
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

func (s *Site) InstallWordPress() error {

	setupCommand := []string{
		"core",
		"install",
		fmt.Sprintf("--url=%s", s.GetURL(false)),
		"--title='Chris Wiegman Theme Development'",
		"--admin_user=admin",
		"--admin_password=password",
		"--admin_email=contact@chriswiegman.com",
	}

	_, err := s.RunWPCli(setupCommand)
	return err

}

func (s *Site) RunWPCli(command []string) (string, error) {

	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return "", err
	}

	siteDir := path.Join(s.appConfig.AppDirectory, "sites", s.appConfig.SiteName)
	appDir := path.Join(siteDir, "app")

	fullCommand := []string{
		"wp",
		"--path=/var/www/html",
	}

	fullCommand = append(fullCommand, command...)

	container := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana_%s_wordpress_cli", s.appConfig.SiteName),
		Image:       "wordpress:cli",
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana_%s_wordpress_cli", s.appConfig.SiteName),
		Command:     fullCommand,
		Env: []string{
			fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.appConfig.SiteName),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
		},
		Labels: map[string]string{
			"kana.site": s.appConfig.SiteName,
		},
		Volumes: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: appDir,
				Target: "/var/www/html",
			},
		},
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
