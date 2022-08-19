package traefik

import (
	"path"

	"github.com/ChrisWiegman/kana/internal/appConfig"
	"github.com/ChrisWiegman/kana/internal/appSetup"
	"github.com/ChrisWiegman/kana/internal/docker"

	"github.com/docker/docker/api/types/mount"
)

var traefikContainerName = "kana_traefik"

type Traefik struct {
	dockerClient docker.DockerClient
	appDirectory string
}

// NewTraefik Setup a new traefik object for controlling the traefik container
func NewTraefik(staticConfig appConfig.StaticConfig) (*Traefik, error) {

	t := new(Traefik)

	err := appSetup.EnsureCerts(staticConfig)
	if err != nil {
		return t, err
	}

	dockerClient, err := docker.NewController()
	if err != nil {
		return t, err
	}

	t.appDirectory = staticConfig.AppDirectory
	t.dockerClient = *dockerClient

	return t, nil
}

// StartTraefik starts the Traefik container
func (t *Traefik) StartTraefik() error {

	_, _, err := t.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	err = t.dockerClient.EnsureImage("traefik")
	if err != nil {
		return err
	}

	traefikPorts := []docker.ExposedPorts{
		{Port: "80", Protocol: "tcp"},
		{Port: "443", Protocol: "tcp"},
		{Port: "8080", Protocol: "tcp"},
	}

	traefikConfig := docker.ContainerConfig{
		Name:        traefikContainerName,
		Image:       "traefik",
		Ports:       traefikPorts,
		NetworkName: "kana",
		HostName:    "kanatraefik",
		Labels: map[string]string{
			"kana.global": "true",
		},
		Volumes: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: path.Join(t.appDirectory, "config", "traefik", "traefik.toml"),
				Target: "/etc/traefik/traefik.toml",
			},
			{
				Type:   mount.TypeBind,
				Source: path.Join(t.appDirectory, "config", "traefik", "dynamic.toml"),
				Target: "/etc/traefik/dynamic.toml",
			},
			{
				Type:   mount.TypeBind,
				Source: path.Join(t.appDirectory, "certs"),
				Target: "/var/certs",
			},
			{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
		},
	}

	_, err = t.dockerClient.ContainerRun(traefikConfig)

	return err
}

// MaybeStopTraefik Checks to see if other sites are running and shuts down the traefik instance if none are
func (t *Traefik) MaybeStopTraefik() error {

	containers, err := t.dockerClient.ListContainers("")
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		return t.StopTraefik()
	}

	return nil
}

// Stops the Traefik container
func (t *Traefik) StopTraefik() error {

	_, err := t.dockerClient.ContainerStop(traefikContainerName)

	return err
}
