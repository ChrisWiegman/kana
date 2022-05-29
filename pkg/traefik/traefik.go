package traefik

import (
	"path"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/pkg/docker"
	"github.com/docker/docker/api/types"
)

func NewTraefik() {

	controller, err := docker.NewController()
	if err != nil {
		panic(err)
	}

	_, _, err = controller.EnsureNetwork("kana")
	if err != nil {
		panic(err)
	}

	err = controller.EnsureImage("traefik")
	if err != nil {
		panic(err)
	}

	err = controller.EnsureImage("tecnativa/docker-socket-proxy")
	if err != nil {
		panic(err)
	}

	traefikPorts := []docker.ExposedPorts{
		{Port: "80", Protocol: "tcp"},
		{Port: "443", Protocol: "tcp"},
	}

	dockerProxyPorts := []docker.ExposedPorts{
		{Port: "2375", Protocol: "tcp"},
		{Port: "8080", Protocol: "tcp"},
	}

	configRoot, _ := config.GetConfigRoot()

	dockerProxyConfig := docker.ContainerConfig{
		Image:       "tecnativa/docker-socket-proxy",
		Ports:       dockerProxyPorts,
		NetworkName: "kana",
		Volumes: []docker.VolumeMount{
			{
				HostPath: "/var/run/docker.sock",
				Volume: &types.Volume{
					Mountpoint: "/var/run/docker.sock",
				},
			},
		},
		Command: []string{},
	}

	traefikConfig := docker.ContainerConfig{
		Image:       "traefik",
		Ports:       traefikPorts,
		NetworkName: "kana",
		Volumes: []docker.VolumeMount{
			{
				HostPath: path.Join(configRoot, ".kana", "conf", "traefik", "traefik.toml"),
				Volume: &types.Volume{
					Mountpoint: "/etc/traefik/traefik.toml",
				},
			},
			{
				HostPath: path.Join(configRoot, ".kana", "conf", "traefik", "dynamic.toml"),
				Volume: &types.Volume{
					Mountpoint: "/etc/traefik/traefik.toml",
				},
			},
		},
		Command: []string{},
	}

	_, err = controller.ContainerRun(dockerProxyConfig)
	if err != nil {
		panic(err)
	}

	_, err = controller.ContainerRun(traefikConfig)
	if err != nil {
		panic(err)
	}
}
