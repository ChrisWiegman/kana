package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
)

func (c *Controller) ContainerRun(image string, ports []PortList, networkName string, volumes []VolumeMount, command []string) (id string, err error) {
	hostConfig := container.HostConfig{}
	containerPorts := c.getNetworkConfig(ports)

	if len(containerPorts.PortBindings) > 0 {
		hostConfig.PortBindings = containerPorts.PortBindings
	}

	networkConfig := network.NetworkingConfig{}

	if len(networkName) > 0 {
		networkConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			networkName: {},
		}
	}

	var mounts []mount.Mount

	for _, volume := range volumes {
		mount := mount.Mount{
			Type:   mount.TypeVolume,
			Source: volume.Volume.Name,
			Target: volume.HostPath,
		}
		mounts = append(mounts, mount)
	}

	hostConfig.Mounts = mounts

	fmt.Println(containerPorts.PortSet)

	resp, err := c.cli.ContainerCreate(context.Background(), &container.Config{
		Tty:          true,
		Image:        image,
		ExposedPorts: containerPorts.PortSet,
		Cmd:          command,
	}, &hostConfig, &networkConfig, nil, "")

	if err != nil {
		return "", err
	}

	err = c.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *Controller) ContainerWait(id string) (state int64, err error) {
	containerResult, errorCode := c.cli.ContainerWait(context.Background(), id, "")
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

	reader, err := c.cli.ContainerLogs(ctx, id, types.ContainerLogsOptions{
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

func (c *Controller) ContainerRunAndClean(image string, ports []PortList, networkName string, volumes []VolumeMount, command []string) (statusCode int64, body string, err error) {

	// Start the container
	id, err := c.ContainerRun(image, ports, networkName, volumes, command)
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

	err = c.cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})

	if err != nil {
		fmt.Printf("Unable to remove container %q: %q\n", id, err)
	}

	return statusCode, body, err
}
