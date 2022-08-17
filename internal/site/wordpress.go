package site

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/traefik"
	"github.com/docker/docker/api/types/mount"
)

func (s *Site) GetSiteContainers() []string {

	return []string{
		fmt.Sprintf("kana_%s_database", s.StaticConfig.SiteName),
		fmt.Sprintf("kana_%s_wordpress", s.StaticConfig.SiteName),
	}
}

func (s *Site) IsSiteRunning() bool {

	containers, _ := s.dockerClient.ListContainers(s.StaticConfig.SiteName)

	return len(containers) != 0

}

func (s *Site) StopWordPress() error {

	wordPressContainers := s.GetSiteContainers()

	for _, wordPressContainer := range wordPressContainers {
		_, err := s.dockerClient.ContainerStop(wordPressContainer)
		if err != nil {
			return err
		}
	}

	traefikClient, err := traefik.NewTraefik(s.StaticConfig)
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

func (s *Site) StartWordPress() error {

	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	appDir := path.Join(s.StaticConfig.SiteDirectory, "app")
	databaseDir := path.Join(s.StaticConfig.SiteDirectory, "database")

	if s.SiteConfig.Local {
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

	if s.SiteConfig.Type == "plugin" {
		appVolumes = append(appVolumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: cwd,
			Target: path.Join("/var/www/html", "wp-content", "plugins", s.StaticConfig.SiteName),
		})
	}

	if s.SiteConfig.Type == "theme" {
		appVolumes = append(appVolumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: cwd,
			Target: path.Join("/var/www/html", "wp-content", "themes", s.StaticConfig.SiteName),
		})
	}

	wordPressContainers := []docker.ContainerConfig{
		{
			Name:        fmt.Sprintf("kana_%s_database", s.StaticConfig.SiteName),
			Image:       "mariadb",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_database", s.StaticConfig.SiteName),
			Env: []string{
				"MARIADB_ROOT_PASSWORD=password",
				"MARIADB_DATABASE=wordpress",
				"MARIADB_USER=wordpress",
				"MARIADB_PASSWORD=wordpress",
			},
			Labels: map[string]string{
				"kana.site": s.StaticConfig.SiteName,
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
			Name:        fmt.Sprintf("kana_%s_wordpress", s.StaticConfig.SiteName),
			Image:       fmt.Sprintf("wordpress:php%s", s.SiteConfig.PHPVersion),
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_wordpress", s.StaticConfig.SiteName),
			Env: []string{
				fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.StaticConfig.SiteName),
				"WORDPRESS_DB_USER=wordpress",
				"WORDPRESS_DB_PASSWORD=wordpress",
				"WORDPRESS_DB_NAME=wordpress",
			},
			Labels: map[string]string{
				"traefik.enable": "true",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", s.StaticConfig.SiteName): "web",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", s.StaticConfig.SiteName):        fmt.Sprintf("Host(`%s.%s`)", s.StaticConfig.SiteName, s.StaticConfig.AppDomain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", s.StaticConfig.SiteName):      "websecure",
				fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", s.StaticConfig.SiteName):             fmt.Sprintf("Host(`%s.%s`)", s.StaticConfig.SiteName, s.StaticConfig.AppDomain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", s.StaticConfig.SiteName):              "true",
				"kana.site": s.StaticConfig.SiteName,
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

	fmt.Println("Finishing WordPress setup...")

	setupCommand := []string{
		"core",
		"install",
		fmt.Sprintf("--url=%s", s.GetURL(false)),
		fmt.Sprintf("--title=Kana Development %s: %s", s.DynamicConfig.SiteType, s.StaticConfig.SiteName),
		fmt.Sprintf("--admin_user=%s", s.DynamicConfig.AdminUsername),
		fmt.Sprintf("--admin_password=%s", s.DynamicConfig.AdminPassword),
		fmt.Sprintf("--admin_email=%s", s.DynamicConfig.AdminEmail),
	}

	_, err := s.RunWPCli(setupCommand)
	return err

}

func (s *Site) RunWPCli(command []string) (string, error) {

	_, _, err := s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return "", err
	}

	siteDir := path.Join(s.StaticConfig.AppDirectory, "sites", s.StaticConfig.SiteName)
	appDir := path.Join(siteDir, "app")

	fullCommand := []string{
		"wp",
		"--path=/var/www/html",
	}

	fullCommand = append(fullCommand, command...)

	container := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana_%s_wordpress_cli", s.StaticConfig.SiteName),
		Image:       fmt.Sprintf("wordpress:cli-php%s", s.DynamicConfig.PHPVersion),
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana_%s_wordpress_cli", s.StaticConfig.SiteName),
		Command:     fullCommand,
		Env: []string{
			fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", s.StaticConfig.SiteName),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
		},
		Labels: map[string]string{
			"kana.site": s.StaticConfig.SiteName,
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
