package site

import (
	"fmt"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/docker"
)

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

// startPHPMyAdmin Starts the PhpMyAdmin container.
func (s *Site) startPHPMyAdmin(consoleOutput *console.Console) error {
	phpMyAdminContainer := s.getPhpMyAdminContainer()

	return s.startContainer(&phpMyAdminContainer, true, false, consoleOutput)
}
