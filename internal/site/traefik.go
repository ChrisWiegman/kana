package site

import (
	"path"

	"github.com/ChrisWiegman/kana-cli/pkg/docker"

	"github.com/docker/docker/api/types/mount"
)

var traefikContainerName = "kana_traefik"

// maybeStopTraefik Checks to see if other sites are running and shuts down the traefik instance if none are
func (s *Site) maybeStopTraefik() error {
	containers, err := s.dockerClient.ListContainers("")
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		return s.stopTraefik()
	}

	return nil
}

// startTraefik Starts the Traefik container
func (s *Site) startTraefik() error {
	err := s.Settings.EnsureSSLCerts()
	if err != nil {
		return err
	}

	_, _, err = s.dockerClient.EnsureNetwork("kana")
	if err != nil {
		return err
	}

	err = s.dockerClient.EnsureImage("traefik")
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
				Source: path.Join(s.Settings.AppDirectory, "config", "traefik", "traefik.toml"),
				Target: "/etc/traefik/traefik.toml",
			},
			{
				Type:   mount.TypeBind,
				Source: path.Join(s.Settings.AppDirectory, "config", "traefik", "dynamic.toml"),
				Target: "/etc/traefik/dynamic.toml",
			},
			{
				Type:   mount.TypeBind,
				Source: path.Join(s.Settings.AppDirectory, "certs"),
				Target: "/var/certs",
			},
			{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
		},
	}

	_, err = s.dockerClient.ContainerRun(&traefikConfig, false, false)

	return err
}

// stopTraefik Stops the Traefik container
func (s *Site) stopTraefik() error {
	_, err := s.dockerClient.ContainerStop(traefikContainerName)
	if err != nil {
		return err
	}

	// Delete the "kana" network as well
	_, err = s.dockerClient.RemoveNetwork("kana")

	return err
}
