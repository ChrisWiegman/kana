package traefik

import (
	"path"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/docker"

	"github.com/docker/docker/api/types/mount"
)

func NewTraefik(controller *docker.Controller) error {

	err := controller.EnsureImage("traefik")
	if err != nil {
		return err
	}

	traefikPorts := []docker.ExposedPorts{
		{Port: "80", Protocol: "tcp"},
		{Port: "443", Protocol: "tcp"},
		{Port: "8080", Protocol: "tcp"},
	}

	configRoot, err := config.GetConfigRoot()
	if err != nil {
		return err
	}

	traefikConfig := docker.ContainerConfig{
		Name:        "kana_traefik",
		Image:       "traefik",
		Ports:       traefikPorts,
		NetworkName: "kana",
		HostName:    "kanatraefik",
		Volumes: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: path.Join(configRoot, "conf", "traefik", "traefik.toml"),
				Target: "/etc/traefik/traefik.toml",
			},
			{
				Type:   mount.TypeBind,
				Source: path.Join(configRoot, "conf", "traefik", "dynamic.toml"),
				Target: "/etc/traefik/dynamic.toml",
			},
			{
				Type:   mount.TypeBind,
				Source: path.Join(configRoot, "certs"),
				Target: "/var/certs",
			},
			{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
		},
	}

	_, err = controller.ContainerRun(traefikConfig)

	return err

}
