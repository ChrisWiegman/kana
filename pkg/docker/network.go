package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

func (c *Controller) EnsureNetwork(name string) (created bool, network *types.NetworkResource, err error) {

	networks, err := c.cli.NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		return false, nil, err
	}

	for _, n := range networks {
		if n.Name == name {
			return true, &n, nil
		}
	}

	networkCreateResults, err := c.cli.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Driver: "bridge",
	})

	if err != nil {
		fmt.Println("fail")
		return false, nil, err
	}

	fmt.Println(networkCreateResults)

	return true, network, nil

}
