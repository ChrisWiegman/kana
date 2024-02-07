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

func (d *Client) ContainerExec(containerName string, rootUser bool, command []string) (ExecResult, error) {
	containerID, isRunning := d.containerIsRunning(containerName)
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

	if rootUser {
		execConfig.User = "root"
	}

	containerResponse, err := d.apiClient.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		return ExecResult{}, err
	}

	execID := containerResponse.ID

	// run it, with stdout/stderr attached
	apiResponse, err := d.apiClient.ContainerExecAttach(context.Background(), execID, types.ExecStartCheck{})
	if err != nil {
		return ExecResult{}, err
	}

	defer apiResponse.Close()

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, apiResponse.Reader)
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
	inspectResponse, err := d.apiClient.ContainerExecInspect(context.Background(), execID)
	if err != nil {
		return ExecResult{}, err
	}

	return ExecResult{
			ExitCode: inspectResponse.ExitCode,
			StdOut:   outBuf.String(),
			StdErr:   errBuf.String(),
		},
		nil
}

// ContainerGetMounts Returns a slice containing all the mounts to the given container.
func (d *Client) ContainerGetMounts(containerName string) []types.MountPoint {
	containerID, isRunning := d.containerIsRunning(containerName)
	if !isRunning {
		return []types.MountPoint{}
	}

	results, _ := d.apiClient.ContainerInspect(context.Background(), containerID)

	return results.Mounts
}

// containerIsRunning Checks if a given container is running by name.
func (d *Client) containerIsRunning(containerName string) (id string, isRunning bool) {
	containers, err := d.apiClient.ContainerList(context.Background(), container.ListOptions{})
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

// ContainerList Lists all running containers for a given site or all sites if no site is specified.
func (d *Client) ContainerList(site string) ([]types.Container, error) {
	f := filters.NewArgs()

	if site == "" {
		f.Add("label", "kana.site")
	} else {
		f.Add("label", fmt.Sprintf("kana.site=%s", site))
	}

	options := container.ListOptions{
		All:     true,
		Filters: f,
	}

	containers, err := d.apiClient.ContainerList(context.Background(), options)

	return containers, err
}

func (d *Client) containerLog(id string) (result string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sleepDuration)*time.Second)
	defer cancel()

	reader, err := d.apiClient.ContainerLogs(ctx, id, container.LogsOptions{
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

func (d *Client) ContainerRestart(containerName string) (bool, error) {
	containerID, isRunning := d.containerIsRunning(containerName)
	if !isRunning {
		return true, nil
	}

	err := d.apiClient.ContainerStop(context.Background(), containerID, container.StopOptions{})
	if err != nil {
		return false, err
	}

	err = d.apiClient.ContainerStart(context.Background(), containerID, container.StartOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *Client) ContainerRun(config *ContainerConfig, randomPorts, localUser bool) (id string, err error) {
	containerID, isRunning := d.containerIsRunning(config.Name)
	if isRunning {
		return containerID, nil
	}

	hostConfig := container.HostConfig{}
	containerPorts, err := getNetworkConfig(config.Ports, randomPorts)
	if err != nil {
		return containerID, err
	}

	if len(containerPorts.PortBindings) > 0 {
		hostConfig.PortBindings = containerPorts.PortBindings
	}

	networkConfig := network.NetworkingConfig{}

	if config.NetworkName != "" {
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

	resp, err := d.apiClient.ContainerCreate(context.Background(), containerConfig, &hostConfig, &networkConfig, nil, config.Name)
	if err != nil {
		return "", err
	}

	err = d.apiClient.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (d *Client) ContainerRunAndClean(config *ContainerConfig) (statusCode int64, body string, err error) {
	// Start the container
	id, err := d.ContainerRun(config, false, true)
	if err != nil {
		return statusCode, body, err
	}

	// Wait for it to finish
	statusCode, err = d.containerWait(id)
	if err != nil {
		return statusCode, body, err
	}

	// Get the log
	body, _ = d.containerLog(id)

	err = d.apiClient.ContainerRemove(context.Background(), id, container.RemoveOptions{})
	return statusCode, body, err
}

func (d *Client) ContainerStop(containerName string) (bool, error) {
	containerID, isRunning := d.containerIsRunning(containerName)
	if !isRunning {
		return true, nil
	}

	err := d.apiClient.ContainerStop(context.Background(), containerID, container.StopOptions{})
	if err != nil {
		return false, err
	}

	err = d.apiClient.ContainerRemove(context.Background(), containerID, container.RemoveOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *Client) containerWait(id string) (state int64, err error) {
	containerResult, errorCode := d.apiClient.ContainerWait(context.Background(), id, "")

	select {
	case err := <-errorCode:
		return 0, err
	case result := <-containerResult:
		return result.StatusCode, nil
	}
}
