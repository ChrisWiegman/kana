package docker

/**
 * Docker code examples currently from https://willschenk.com/articles/2021/controlling_docker_in_golang/
 **/

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var execCommand = exec.Command

var maxRetries = 12
var sleepDuration = 5

// DockerClient is an interface the must be implemented to provide Docker services through this package.
type DockerClient struct {
	moby APIClient
}

func NewDockerClient(consoleOutput *console.Console) (dockerClient *DockerClient, err error) {
	dockerClient = new(DockerClient)

	dockerClient.moby, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	err = ensureDockerIsAvailable(consoleOutput, dockerClient.moby)
	if err != nil {
		return nil, err
	}

	return dockerClient, nil
}

func ensureDockerIsAvailable(consoleOutput *console.Console, moby APIClient) error {
	_, err := moby.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		fmt.Println("We have an error")
		if runtime.GOOS == "darwin" {
			consoleOutput.Println("Docker doesn't appear to be running. Trying to start Docker.")
			err = execCommand("open", "-a", "Docker").Run()
			fmt.Println(err)
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

				_, err = moby.ContainerList(context.Background(), types.ContainerListOptions{})
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}
