package wordpress

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/traefik"
	"github.com/docker/docker/api/types/mount"
)

func (s *KanaSite) StopWordPress() error {

	wordPressContainers := []string{
		fmt.Sprintf("kana_%s_database", s.controller.Config.SiteDirectory),
		fmt.Sprintf("kana_%s_wordpress", s.controller.Config.SiteDirectory),
	}

	for _, wordPressContainer := range wordPressContainers {
		_, err := s.controller.ContainerStop(wordPressContainer)
		if err != nil {
			return err
		}
	}

	return traefik.MaybeStopTraefik(s.controller)

}

func (s *KanaSite) StartWordPress() error {

	_, _, err := s.controller.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	siteDir := path.Join(s.controller.Config.AppDirectory, "sites", s.controller.Config.SiteDirectory)
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
			Name:        fmt.Sprintf("kana_%s_database", s.controller.Config.SiteDirectory),
			Image:       "mariadb",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_database", s.controller.Config.SiteDirectory),
			Env: []string{
				"MARIADB_ROOT_PASSWORD=password",
				"MARIADB_DATABASE=wordpress",
				"MARIADB_USER=wordpress",
				"MARIADB_PASSWORD=wordpress",
			},
			Labels: map[string]string{
				"kana.site": s.controller.Config.SiteDirectory,
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
			Name:        fmt.Sprintf("kana_%s_wordpress", s.controller.Config.SiteDirectory),
			Image:       "wordpress",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_wordpress", s.controller.Config.SiteDirectory),
			Env: []string{
				fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.controller.Config.SiteDirectory),
				"WORDPRESS_DB_USER=wordpress",
				"WORDPRESS_DB_PASSWORD=wordpress",
				"WORDPRESS_DB_NAME=wordpress",
			},
			Labels: map[string]string{
				"traefik.enable": "true",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", s.controller.Config.SiteDirectory): "web",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", s.controller.Config.SiteDirectory):        fmt.Sprintf("Host(`%s.%s`)", s.controller.Config.SiteDirectory, s.controller.Config.AppDomain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", s.controller.Config.SiteDirectory):      "websecure",
				fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", s.controller.Config.SiteDirectory):             fmt.Sprintf("Host(`%s.%s`)", s.controller.Config.SiteDirectory, s.controller.Config.AppDomain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", s.controller.Config.SiteDirectory):              "true",
				"kana.site": s.controller.Config.SiteDirectory,
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

		err := s.controller.EnsureImage(container.Image)
		if err != nil {
			return err
		}

		_, err = s.controller.ContainerRun(container)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *KanaSite) InstallWordPress() error {

	setupCommand := []string{
		"core",
		"install",
		fmt.Sprintf("--url=%s", s.GetURL(false)),
		"--title='Chris Wiegman Theme Development'",
		"--admin_user=admin",
		"--admin_password=password",
		"--admin_email=contact@chriswiegman.com",
	}

	return s.RunCli(setupCommand)

}

func (s *KanaSite) RunCli(command []string) error {

	_, _, err := s.controller.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	siteDir := path.Join(s.controller.Config.AppDirectory, "sites", s.controller.Config.SiteDirectory)
	appDir := path.Join(siteDir, "app")

	fullCommand := []string{
		"wp",
		"--path=/var/www/html",
	}

	fullCommand = append(fullCommand, command...)

	container := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana_%s_wordpress_cli", s.controller.Config.SiteDirectory),
		Image:       "wordpress:cli",
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana_%s_wordpress_cli", s.controller.Config.SiteDirectory),
		Command:     fullCommand,
		Env: []string{
			fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.controller.Config.SiteDirectory),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
		},
		Labels: map[string]string{
			"kana.site": s.controller.Config.SiteDirectory,
		},
		Volumes: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: appDir,
				Target: "/var/www/html",
			},
		},
	}

	err = s.controller.EnsureImage(container.Image)
	if err != nil {
		return err
	}

	_, output, err := s.controller.ContainerRunAndClean(container)
	if err != nil {
		return err
	}

	fmt.Println(output)

	return nil

}
