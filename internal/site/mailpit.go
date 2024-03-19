package site

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/docker"

	"github.com/docker/docker/api/types/mount"
)

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
			"kana.type":      "mailpit",
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

func (s *Site) isMailpitRunning() bool {
	// We need container details to see if the mailpit container is running
	containers, err := s.dockerClient.ContainerList(s.Settings.Name)
	if err != nil {
		return false
	}

	for i := range containers {
		if containers[i].Image == "axllent/mailpit" {
			return true
		}
	}

	return false
}

// startMailpit Starts the Mailpit container.
func (s *Site) startMailpit(consoleOutput *console.Console) error {
	mailpitContainer := s.getMailpitContainer()

	return s.startContainer(&mailpitContainer, true, true, consoleOutput)
}
