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

	traefikPorts := []docker.PortList{
		{Port: "80", Protocol: "tcp"},
		{Port: "443", Protocol: "tcp"},
	}

	_, err = controller.ContainerRun("traefik", traefikPorts, "kana", []docker.VolumeMount{}, []string{})
	if err != nil {
		panic(err)
	}

}
