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
		fmt.Sprintf("kana_%s_database", s.appConfig.SiteDirectory),
		fmt.Sprintf("kana_%s_wordpress", s.appConfig.SiteDirectory),
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

func (s *Site) StartWordPress() error {

	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	siteDir := path.Join(s.appConfig.AppDirectory, "sites", s.appConfig.SiteDirectory)
	appDir := path.Join(siteDir, "app")
	databaseDir := path.Join(siteDir, "database")

	if err := os.MkdirAll(appDir, 0750); err != nil {
		return err
	}

	if err := os.MkdirAll(databaseDir, 0750); err != nil {
		return err
	}

	wordPressContainers := []docker.ContainerConfig{
		{
			Name:        fmt.Sprintf("kana_%s_database", s.appConfig.SiteDirectory),
			Image:       "mariadb",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_database", s.appConfig.SiteDirectory),
			Env: []string{
				"MARIADB_ROOT_PASSWORD=password",
				"MARIADB_DATABASE=wordpress",
				"MARIADB_USER=wordpress",
				"MARIADB_PASSWORD=wordpress",
			},
			Labels: map[string]string{
				"kana.site": s.appConfig.SiteDirectory,
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
			Name:        fmt.Sprintf("kana_%s_wordpress", s.appConfig.SiteDirectory),
			Image:       "wordpress",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_wordpress", s.appConfig.SiteDirectory),
			Env: []string{
				fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.appConfig.SiteDirectory),
				"WORDPRESS_DB_USER=wordpress",
				"WORDPRESS_DB_PASSWORD=wordpress",
				"WORDPRESS_DB_NAME=wordpress",
			},
			Labels: map[string]string{
				"traefik.enable": "true",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", s.appConfig.SiteDirectory): "web",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", s.appConfig.SiteDirectory):        fmt.Sprintf("Host(`%s.%s`)", s.appConfig.SiteDirectory, s.appConfig.AppDomain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", s.appConfig.SiteDirectory):      "websecure",
				fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", s.appConfig.SiteDirectory):             fmt.Sprintf("Host(`%s.%s`)", s.appConfig.SiteDirectory, s.appConfig.AppDomain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", s.appConfig.SiteDirectory):              "true",
				"kana.site": s.appConfig.SiteDirectory,
			},
			Volumes: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: appDir,
					Target: "/var/www/html",
				},
			},
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

	return s.RunWPCli(setupCommand)

}

func (s *Site) RunWPCli(command []string) error {

	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	siteDir := path.Join(s.appConfig.AppDirectory, "sites", s.appConfig.SiteDirectory)
	appDir := path.Join(siteDir, "app")

	fullCommand := []string{
		"wp",
		"--path=/var/www/html",
	}

	fullCommand = append(fullCommand, command...)

	container := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana_%s_wordpress_cli", s.appConfig.SiteDirectory),
		Image:       "wordpress:cli",
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana_%s_wordpress_cli", s.appConfig.SiteDirectory),
		Command:     fullCommand,
		Env: []string{
			fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.appConfig.SiteDirectory),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
		},
		Labels: map[string]string{
			"kana.site": s.appConfig.SiteDirectory,
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
		return err
	}

	_, output, err := s.dockerClient.ContainerRunAndClean(container)
	if err != nil {
		return err
	}

	fmt.Println(output)

	return nil

}
