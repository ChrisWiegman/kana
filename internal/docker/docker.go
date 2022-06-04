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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Controller struct {
	cli *client.Client
}

func NewController() (c *Controller, err error) {
	c = new(Controller)

	c.cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Here is the error")
		return nil, err
	}

	ensureDockerIsAvailable(c)

	return c, nil
}

func ensureDockerIsAvailable(c *Controller) {

	_, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		if runtime.GOOS == "darwin" {

			fmt.Println("Docker doesn't appear to be running. Trying to start Docker.")
			exec.Command("open", "-a", "Docker").Run()

			retries := 0

			for retries <= 12 {

				retries++

				if retries == 12 {
					fmt.Println("Restarting Docker is taking too long. We seem to have hit an error")
					return
				}

				time.Sleep(5 * time.Second)

				_, err = c.cli.ContainerList(context.Background(), types.ContainerListOptions{})
				if err == nil {
					break
				}
			}
		}
	}
}
