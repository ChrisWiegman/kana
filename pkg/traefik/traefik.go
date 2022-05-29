package traefik

import (
	"path"

	"github.com/ChrisWiegman/kana/internal/templates"
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

	traefikPorts := []docker.ExposedPorts{
		{Port: "80", Protocol: "tcp"},
		{Port: "443", Protocol: "tcp"},
	}

	configRoot, _ := templates.GetConfigRoot()

	config := docker.ContainerConfig{
		Image:       "traefik",
		Ports:       traefikPorts,
		NetworkName: "kana",
		Volumes: []docker.VolumeMount{
			{
				HostPath: path.Join(configRoot, ".kana", "conf", "traefik", "traefik.yaml"),
				Volume: &types.Volume{
					Mountpoint: "/etc/traefik/traefik.yml",
				},
			},
		},
		Command: []string{},
	}

	_, err = controller.ContainerRun(config)
	if err != nil {
		panic(err)
	}
}
