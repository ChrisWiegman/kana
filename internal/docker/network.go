package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
)

type ExposedPorts struct {
	Port     string
	Protocol string
}

type portConfig struct {
	PortBindings nat.PortMap
	PortSet      nat.PortSet
}

func (d *DockerClient) EnsureNetwork(name string) (created bool, network types.NetworkResource, err error) {
	hasNetwork, network, err := findNetworkByName(name, d.moby)

	if err != nil {
		return false, types.NetworkResource{}, err
	}

	if hasNetwork {
		return false, network, nil
	}

	networkCreateResults, err := d.moby.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Driver: "bridge",
	})

	if err != nil {
		return false, types.NetworkResource{}, err
	}

	hasNetwork, network, err = findNetworkByID(networkCreateResults.ID, d.moby)

	if err != nil {
		return false, types.NetworkResource{}, err
	}

	if hasNetwork {
		return true, network, nil
	}

	return false, types.NetworkResource{}, fmt.Errorf("could not create network")
}

func (d *DockerClient) RemoveNetwork(name string) (removed bool, err error) {
	hasNetwork, network, err := findNetworkByName(name, d.moby)

	if err != nil {
		return false, err
	}

	if !hasNetwork {
		return false, nil
	}

	return true, d.moby.NetworkRemove(context.Background(), network.ID)
}

func findNetworkByID(id string, moby APIClient) (found bool, network types.NetworkResource, err error) {
	networks, err := moby.NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		return false, types.NetworkResource{}, err
	}

	for i := range networks {
		if networks[i].ID == id {
			return true, networks[i], nil
		}
	}

	return false, types.NetworkResource{}, nil
}

func findNetworkByName(name string, moby APIClient) (found bool, network types.NetworkResource, err error) {
	networks, err := moby.NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		return false, types.NetworkResource{}, err
	}

	for i := range networks {
		if networks[i].Name == name {
			return true, networks[i], nil
		}
	}

	return false, types.NetworkResource{}, nil
}

func getNetworkConfig(ports []ExposedPorts, randomPorts bool) portConfig {
	portBindings := make(nat.PortMap)
	portSet := make(nat.PortSet)

	for _, port := range ports {
		portName, err := nat.NewPort(port.Protocol, port.Port)
		if err != nil {
			panic(err)
		}

		hostPort := port.Port

		// Reassign host port to "0" to chose a random port where applicable
		if randomPorts {
			hostPort = "0"
		}

		portBindings[portName] = []nat.PortBinding{
			{
				HostPort: hostPort,
			},
		}

		portSet[portName] = struct{}{}
	}

	return portConfig{
		PortBindings: portBindings,
		PortSet:      portSet,
	}
}