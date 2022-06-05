package docker

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
)

type ContainerConfig struct {
	Name        string
	Image       string
	Ports       []ExposedPorts
	HostName    string
	NetworkName string
	Volumes     []mount.Mount
	Command     []string
	Env         []string
	Labels      map[string]string
}

func (c *Controller) IsContainerRunning(containerName string) (id string, isRunning bool) {

	containers, err := c.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return "", false
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if containerName == strings.Trim(name, "/") {
				return container.ID, true
			}
		}
	}

	return "", false

}

func (c *Controller) ContainerRun(config ContainerConfig) (id string, err error) {

	containerID, isRunning := c.IsContainerRunning(config.Name)
	if isRunning {
		return containerID, nil
	}

	hostConfig := container.HostConfig{}
	containerPorts := c.getNetworkConfig(config.Ports)

	if len(containerPorts.PortBindings) > 0 {
		hostConfig.PortBindings = containerPorts.PortBindings
	}

	networkConfig := network.NetworkingConfig{}

	if len(config.NetworkName) > 0 {
		networkConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			config.NetworkName: {},
		}
	}

	hostConfig.Mounts = config.Volumes

	resp, err := c.client.ContainerCreate(context.Background(), &container.Config{
		Tty:          true,
		Image:        config.Image,
		ExposedPorts: containerPorts.PortSet,
		Cmd:          config.Command,
		Hostname:     config.HostName,
		Env:          config.Env,
		Labels:       config.Labels,
	}, &hostConfig, &networkConfig, nil, config.Name)

	if err != nil {
		return "", err
	}

	err = c.client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *Controller) ContainerWait(id string) (state int64, err error) {
	containerResult, errorCode := c.client.ContainerWait(context.Background(), id, "")
	select {
	case err := <-errorCode:
		return 0, err
	case result := <-containerResult:
		return result.StatusCode, nil
	}
}

func (c *Controller) ContainerLog(id string) (result string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader, err := c.client.ContainerLogs(ctx, id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true})

	if err != nil {
		return "", err
	}

	buffer, err := io.ReadAll(reader)

	if err != nil && err != io.EOF {
		return "", err
	}

	return string(buffer), nil
}

func (c *Controller) ContainerRunAndClean(config ContainerConfig) (statusCode int64, body string, err error) {

	// Start the container
	id, err := c.ContainerRun(config)
	if err != nil {
		return statusCode, body, err
	}

	// Wait for it to finish
	statusCode, err = c.ContainerWait(id)
	if err != nil {
		return statusCode, body, err
	}

	// Get the log
	body, _ = c.ContainerLog(id)

	err = c.client.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})

	if err != nil {
		fmt.Printf("Unable to remove container %q: %q\n", id, err)
	}

	return statusCode, body, err
}

func (c *Controller) ContainerStop(containerName string) (bool, error) {

	containerID, isRunning := c.IsContainerRunning(containerName)
	if !isRunning {
		return true, nil
	}

	err := c.client.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		return false, err
	}

	err = c.client.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{})
	if err != nil {
		return false, err
	}

	return true, nil

}
