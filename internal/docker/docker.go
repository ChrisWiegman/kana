package docker

/**
 * Docker code examples currently from https://willschenk.com/articles/2021/controlling_docker_in_golang/
 **/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/ChrisWiegman/kana-cli/internal/console"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var execCommand = exec.Command

var maxRetries = 12
var sleepDuration = 5

// DockerClient is an interface the must be implemented to provide Docker services through this package.
type DockerClient struct {
	apiClient       APIClient
	imageUpdateData ViperClient
	checkedImages   []string
}

type DockerContext struct {
	Current        bool   `json:"Current"`
	DockerEndpoint string `json:"DockerEndpoint"`
}

func NewDockerClient(consoleOutput *console.Console, appDirectory string) (dockerClient *DockerClient, err error) {
	dockerClient = new(DockerClient)

	var dockerEndpoint string
	useDefaultHost := false

	dockerEndpoint, err = getCurrentDockerEndpoint()
	if err != nil {
		if err.Error() == "docker context was not found" {
			useDefaultHost = true
		} else {
			return nil, err
		}
	}

	if useDefaultHost {
		dockerClient.apiClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	} else {
		dockerClient.apiClient, err = client.NewClientWithOpts(client.WithHost(dockerEndpoint), client.WithAPIVersionNegotiation())
	}

	if err != nil {
		return nil, err
	}

	err = ensureDockerIsAvailable(consoleOutput, dockerClient.apiClient)
	if err != nil {
		return nil, err
	}

	dockerClient.imageUpdateData, _ = dockerClient.loadImageUpdateData(appDirectory)

	return dockerClient, nil
}

func getCurrentDockerEndpoint() (string, error) {
	dockerContexts := execCommand(
		"docker",
		"context",
		"ls",
		"--format",
		"json")

	var out bytes.Buffer
	dockerContexts.Stdout = &out

	err := dockerContexts.Run()
	if err != nil {
		return "", err
	}

	var contexts []DockerContext

	err = json.Unmarshal(out.Bytes(), &contexts)
	if err != nil {
		var singleContext DockerContext

		err = json.Unmarshal(out.Bytes(), &singleContext)
		if err != nil {
			return "", err
		}

		return singleContext.DockerEndpoint, nil
	}

	for _, context := range contexts {
		if context.Current {
			return context.DockerEndpoint, nil
		}
	}

	return "", fmt.Errorf("docker context was not found")
}

func ensureDockerIsAvailable(consoleOutput *console.Console, apiClient APIClient) error {
	_, err := apiClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		if runtime.GOOS == "darwin" { //nolint:goconst
			consoleOutput.Println("Docker doesn't appear to be running. Trying to start Docker.")
			err = execCommand("open", "-a", "Docker").Run()
			if err != nil {
				return fmt.Errorf("error: unable to start Docker for Mac")
			}

			retries := 0

			for retries <= maxRetries {
				retries++

				if retries == maxRetries {
					consoleOutput.Println("Restarting Docker is taking too long. We seem to have hit an error")
					return fmt.Errorf("error: unable to start Docker for Mac")
				}

				time.Sleep(time.Duration(sleepDuration) * time.Second)

				_, err = apiClient.ContainerList(context.Background(), types.ContainerListOptions{})
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}
