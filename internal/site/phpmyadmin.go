package site

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/docker"
)

func (s *Site) getPhpMyAdminContainer() docker.ContainerConfig {
	phpMyAdminContainer := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana-%s-phpmyadmin", s.settings.Get("Name")),
		Image:       "phpmyadmin",
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana-%s-phpmyadmin", s.settings.Get("Name")),
		Env: []string{
			"MYSQL_ROOT_PASSWORD=password",
			fmt.Sprintf("PMA_HOST=kana-%s-database", s.settings.Get("Name")),
			"PMA_USER=wordpress",
			"PMA_PASSWORD=wordpress",
		},
		Labels: map[string]string{
			"traefik.enable": "true",
			"kana.type":      "phpmyadmin",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-%s-http.entrypoints", s.settings.Get("Name"), "phpmyadmin"): "web",
			fmt.Sprintf(
				"traefik.http.routers.wordpress-%s-%s-http.rule",
				s.settings.Get("Name"),
				"phpmyadmin"): fmt.Sprintf(
				"Host(`%s-%s`)",
				"phpmyadmin",
				s.settings.GetDomain()),
			fmt.Sprintf("traefik.http.routers.wordpress-%s-%s.entrypoints", s.settings.Get("Name"), "phpmyadmin"): "websecure",
			fmt.Sprintf(
				"traefik.http.routers.wordpress-%s-%s.rule",
				s.settings.Get("Name"),
				"phpmyadmin"): fmt.Sprintf(
				"Host(`%s-%s`)",
				"phpmyadmin",
				s.settings.GetDomain()),
			fmt.Sprintf("traefik.http.routers.wordpress-%s-%s.tls", s.settings.Get("Name"), "phpmyadmin"): "true",
			"kana.site": s.settings.Get("Name"),
		},
	}

	return phpMyAdminContainer
}

// startPHPMyAdmin Starts the PhpMyAdmin container.
func (s *Site) startPHPMyAdmin(consoleOutput *console.Console) error {
	phpMyAdminContainer := s.getPhpMyAdminContainer()

	return s.startContainer(&phpMyAdminContainer, true, false, consoleOutput)
}
