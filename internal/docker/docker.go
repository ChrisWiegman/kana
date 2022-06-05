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

	"github.com/ChrisWiegman/kana/internal/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Controller struct {
	client *client.Client
	Config config.KanaConfig
}

func NewController(kanaConfig config.KanaConfig) (c *Controller, err error) {

	c = new(Controller)

	c.Config = kanaConfig

	c.client, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	err = ensureDockerIsAvailable(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func ensureDockerIsAvailable(c *Controller) error {

	_, err := c.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		if runtime.GOOS == "darwin" {

			fmt.Println("Docker doesn't appear to be running. Trying to start Docker.")
			err = exec.Command("open", "-a", "Docker").Run()
			if err != nil {
				return fmt.Errorf("error: unable to start Docker for Mac")
			}

			retries := 0

			for retries <= 12 {

				retries++

				if retries == 12 {
					fmt.Println("Restarting Docker is taking too long. We seem to have hit an error")
					return fmt.Errorf("error: unable to start Docker for Mac")
				}

				time.Sleep(5 * time.Second)

				_, err = c.client.ContainerList(context.Background(), types.ContainerListOptions{})
				if err == nil {
					return err
				}
			}
		}
	}

	return err

}
