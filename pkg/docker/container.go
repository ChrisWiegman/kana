package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/pkg/stdcopy"
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

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

// ListContainers Lists all running containers for a given site or all sites if no site is specified
func (d *DockerClient) ListContainers(site string) ([]types.Container, error) {
	f := filters.NewArgs()

	if site == "" {
		f.Add("label", "kana.site")
	} else {
		f.Add("label", fmt.Sprintf("kana.site=%s", site))
	}

	options := types.ContainerListOptions{
		All:     true,
		Filters: f,
	}

	containers, err := d.client.ContainerList(context.Background(), options)

	return containers, err
}

// IsContainerRunning Checks if a given container is running by name
func (d *DockerClient) IsContainerRunning(containerName string) (id string, isRunning bool) {
	containers, err := d.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return "", false
	}

	for i := range containers {
		for _, name := range containers[i].Names {
			if containerName == strings.Trim(name, "/") {
				return containers[i].ID, true
			}
		}
	}

	return "", false
}

// ContainerGetMounts Returns a slice containing all the mounts to the given container
func (d *DockerClient) ContainerGetMounts(containerName string) []types.MountPoint {
	containerID, isRunning := d.IsContainerRunning(containerName)
	if !isRunning {
		return []types.MountPoint{}
	}

	results, _ := d.client.ContainerInspect(context.Background(), containerID)

	return results.Mounts
}

func (d *DockerClient) ContainerRun(config *ContainerConfig, randomPorts, localUser bool) (id string, err error) {
	containerID, isRunning := d.IsContainerRunning(config.Name)
	if isRunning {
		return containerID, nil
	}

	hostConfig := container.HostConfig{}
	containerPorts := d.getNetworkConfig(config.Ports, randomPorts)

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

	containerConfig := &container.Config{
		Tty:          true,
		Image:        config.Image,
		ExposedPorts: containerPorts.PortSet,
		Cmd:          config.Command,
		Hostname:     config.HostName,
		Env:          config.Env,
		Labels:       config.Labels,
	}

	// Linux doesn't abstract the user so we have to do it ourselves
	if localUser && runtime.GOOS == "linux" {
		var currentUser *user.User

		currentUser, err = user.Current()
		if err != nil {
			return containerID, err
		}

		containerConfig.User = fmt.Sprintf("%s:%s", currentUser.Uid, currentUser.Gid)
	}

	resp, err := d.client.ContainerCreate(context.Background(), containerConfig, &hostConfig, &networkConfig, nil, config.Name)
	if err != nil {
		return "", err
	}

	err = d.client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (d *DockerClient) ContainerWait(id string) (state int64, err error) {
	containerResult, errorCode := d.client.ContainerWait(context.Background(), id, "")

	select {
	case err := <-errorCode:
		return 0, err
	case result := <-containerResult:
		return result.StatusCode, nil
	}
}

func (d *DockerClient) ContainerLog(id string) (result string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sleepDuration)*time.Second)
	defer cancel()

	reader, err := d.client.ContainerLogs(ctx, id, types.ContainerLogsOptions{
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

func (d *DockerClient) ContainerRunAndClean(config *ContainerConfig) (statusCode int64, body string, err error) {
	// Start the container
	id, err := d.ContainerRun(config, false, true)
	if err != nil {
		return statusCode, body, err
	}

	// Wait for it to finish
	statusCode, err = d.ContainerWait(id)
	if err != nil {
		return statusCode, body, err
	}

	// Get the log
	body, _ = d.ContainerLog(id)

	err = d.client.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})

	if err != nil {
		fmt.Printf("Unable to remove container %q: %q\n", id, err)
	}

	return statusCode, body, err
}

func (d *DockerClient) ContainerStop(containerName string) (bool, error) {
	containerID, isRunning := d.IsContainerRunning(containerName)
	if !isRunning {
		return true, nil
	}

	err := d.client.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		return false, err
	}

	err = d.client.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *DockerClient) ContainerRestart(containerName string) (bool, error) {
	containerID, isRunning := d.IsContainerRunning(containerName)
	if !isRunning {
		return true, nil
	}

	err := d.client.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		return false, err
	}

	err = d.client.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *DockerClient) ContainerExec(containerName string, command []string) (ExecResult, error) {
	containerID, isRunning := d.IsContainerRunning(containerName)
	if !isRunning {
		return ExecResult{}, nil
	}

	fullCommand := []string{
		"sh",
		"-c",
	}

	fullCommand = append(fullCommand, command...)

	// prepare exec
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          strslice.StrSlice(fullCommand),
	}

	cresp, err := d.client.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		return ExecResult{}, err
	}

	execID := cresp.ID

	// run it, with stdout/stderr attached
	aresp, err := d.client.ContainerExecAttach(context.Background(), execID, types.ExecStartCheck{})
	if err != nil {
		return ExecResult{}, err
	}

	defer aresp.Close()

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, aresp.Reader)
		outputDone <- err
	}()

	select {
	case err = <-outputDone:
		if err != nil {
			return ExecResult{}, err
		}
		break

	case <-context.Background().Done():
		return ExecResult{}, context.Background().Err()
	}

	// get the exit code
	iresp, err := d.client.ContainerExecInspect(context.Background(), execID)
	if err != nil {
		return ExecResult{}, err
	}

	return ExecResult{
			ExitCode: iresp.ExitCode,
			StdOut:   outBuf.String(),
			StdErr:   errBuf.String(),
		},
		nil
}
