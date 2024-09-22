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
	"strings"

	"github.com/ChrisWiegman/kana-wp/internal/console"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/knadh/koanf/v2"
)

var execCommand = exec.Command

var sleepDuration = 5

// Client is an interface the must be implemented to provide Docker services through this package.
type Client struct {
	apiClient       APIClient
	imageUpdateData *koanf.Koanf
	checkedImages   []string
}

type Context struct {
	Current        bool   `json:"Current"`
	DockerEndpoint string `json:"DockerEndpoint"`
}

func New(consoleOutput *console.Console, appDirectory string) (dockerClient *Client, err error) {
	dockerClient = new(Client)

	var dockerEndpoint string

	dockerEndpoint, err = getCurrentDockerEndpoint()
	if err != nil && err.Error() != "docker context was not found. using default" {
		return nil, err
	}

	dockerClient.apiClient, err = client.NewClientWithOpts(client.WithHost(dockerEndpoint), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	err = ensureDockerIsAvailable(dockerClient.apiClient)
	if err != nil {
		return nil, err
	}

	dockerClient.imageUpdateData, _ = dockerClient.loadImageUpdateData(appDirectory)

	return dockerClient, nil
}

func getCurrentDockerEndpoint() (string, error) {
	rawDockerContexts := execCommand(
		"docker",
		"context",
		"ls",
		"--format",
		"json")

	var out bytes.Buffer
	rawDockerContexts.Stdout = &out

	err := rawDockerContexts.Run()
	if err != nil {
		return client.DefaultDockerHost, err
	}

	var contexts []Context

	err = json.Unmarshal(out.Bytes(), &contexts)
	if err != nil {
		// Docker Desktop and Docker Engine treat the json output differently so we need to account for each.
		rawContexts := strings.Split(out.String(), "\n")

		for i := range rawContexts {
			if rawContexts[i] != "" { // Don't test an empty context
				var singleContext Context

				err = json.Unmarshal([]byte(rawContexts[i]), &singleContext)
				if err != nil {
					return client.DefaultDockerHost, err
				}

				contexts = append(contexts, singleContext)
			}
		}
	}

	for i := range contexts {
		if contexts[i].Current {
			return contexts[i].DockerEndpoint, nil
		}
	}

	return client.DefaultDockerHost, fmt.Errorf("docker context was not found. using default")
}

func ensureDockerIsAvailable(apiClient APIClient) error {
	_, err := apiClient.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return fmt.Errorf("Could not connect to Docker. Is Docker running?") //nolint:stylecheck
	}

	return nil
}
