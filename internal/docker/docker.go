package docker

/**
 * Docker code examples currently from https://willschenk.com/articles/2021/controlling_docker_in_golang/
 **/

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"time"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	client *client.Client
}

func NewController() (c *DockerClient, err error) {

	c = new(DockerClient)

	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	// Docker Desktop 4.13 removes /var/run/docker.sock. This workaround should fix the problem. See https://docs.docker.com/desktop/release-notes/#docker-desktop-4130
	_, err = os.Stat("/var/run/docker.sock")
	if err != nil && os.IsNotExist(err) {
		dockerHost := fmt.Sprintf("unix:///Users/%s/.docker/run/docker.sock", currentUser.Username)
		os.Setenv("DOCKER_HOST", dockerHost)
	}

	c.client, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	err = c.ensureDockerIsAvailable()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (d *DockerClient) ensureDockerIsAvailable() error {

	_, err := d.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		if runtime.GOOS == "darwin" {

			console.Println("Docker doesn't appear to be running. Trying to start Docker.")
			err = exec.Command("open", "-a", "Docker").Run()
			if err != nil {
				return fmt.Errorf("error: unable to start Docker for Mac")
			}

			retries := 0

			for retries <= 12 {

				retries++

				if retries == 12 {
					console.Println("Restarting Docker is taking too long. We seem to have hit an error")
					return fmt.Errorf("error: unable to start Docker for Mac")
				}

				time.Sleep(5 * time.Second)

				_, err = d.client.ContainerList(context.Background(), types.ContainerListOptions{})
				if err == nil {
					return err
				}
			}
		}
	}

	return err
}
