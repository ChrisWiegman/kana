package traefik

import "github.com/ChrisWiegman/kana/pkg/docker"

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

	config := docker.ContainerConfig{
		Image:       "traefik",
		Ports:       traefikPorts,
		NetworkName: "kana",
		Volumes:     []docker.VolumeMount{},
		Command:     []string{},
	}

	_, err = controller.ContainerRun(config)
	if err != nil {
		panic(err)
	}
}
