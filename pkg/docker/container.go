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
	"github.com/docker/go-connections/nat"
)

func (c *Controller) ContainerRun(image string, command []string, volumes []VolumeMount) (id string, err error) {
	hostConfig := container.HostConfig{
		PortBindings: nat.PortMap{
			"80/tcp": {
				nat.PortBinding{
					HostIP:   "127.0.0.1",
					HostPort: "80",
				},
			},
			"443/tcp": {
				nat.PortBinding{
					HostIP:   "127.0.0.1",
					HostPort: "443",
				},
			},
		},
	}
	name := "kana"
	ports := make(nat.PortSet)
	networkConfig := network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			name: {},
		},
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

	resp, err := c.cli.ContainerCreate(context.Background(), &container.Config{
		Tty:          true,
		Image:        image,
		Cmd:          command,
		ExposedPorts: ports,
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

func (c *Controller) ContainerRunAndClean(image string, command []string, volumes []VolumeMount) (statusCode int64, body string, err error) {

	// Start the container
	id, err := c.ContainerRun(image, command, volumes)
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
