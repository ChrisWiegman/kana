package docker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/docker/docker/api/types/network"
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

func (d *Client) EnsureNetwork(name string) (created bool, dockerNetwork network.Inspect, err error) {
	hasNetwork, dockerNetwork, err := findNetworkByName(name, d.apiClient)

	if err != nil {
		return false, network.Inspect{}, err
	}

	if hasNetwork {
		return false, dockerNetwork, nil
	}

	networkCreateResults, err := d.apiClient.NetworkCreate(context.Background(), name, network.CreateOptions{
		Driver: "bridge",
	})

	if err != nil {
		return false, network.Inspect{}, err
	}

	hasNetwork, dockerNetwork, err = findNetworkByID(networkCreateResults.ID, d.apiClient)

	if err != nil {
		return false, network.Inspect{}, err
	}

	if hasNetwork {
		return true, dockerNetwork, nil
	}

	return false, network.Inspect{}, fmt.Errorf("could not create network")
}

func (d *Client) RemoveNetwork(name string) (removed bool, err error) {
	hasNetwork, dockerNetwork, err := findNetworkByName(name, d.apiClient)

	if err != nil {
		return false, err
	}

	if !hasNetwork {
		return false, nil
	}

	return true, d.apiClient.NetworkRemove(context.Background(), dockerNetwork.ID)
}

func findNetworkByID(id string, apiClient APIClient) (found bool, dockerNetwork network.Inspect, err error) {
	dockerNetworks, err := apiClient.NetworkList(context.Background(), network.ListOptions{})

	if err != nil {
		return false, network.Inspect{}, err
	}

	for i := range dockerNetworks {
		if dockerNetworks[i].ID == id {
			return true, dockerNetworks[i], nil
		}
	}

	return false, network.Inspect{}, nil
}

func findNetworkByName(name string, apiClient APIClient) (found bool, dockerNetwork network.Inspect, err error) {
	networks, err := apiClient.NetworkList(context.Background(), network.ListOptions{})

	if err != nil {
		return false, network.Inspect{}, err
	}

	for i := range networks {
		if networks[i].Name == name {
			return true, networks[i], nil
		}
	}

	return false, network.Inspect{}, nil
}

func getNetworkConfig(ports []ExposedPorts, randomPorts bool) (portConfig, error) {
	portBindings := make(nat.PortMap)
	portSet := make(nat.PortSet)

	for _, port := range ports {
		portName, err := nat.NewPort(port.Protocol, port.Port)
		if err != nil {
			return portConfig{}, err
		}

		hostPort := port.Port

		if randomPorts {
			port, err := getRandomPort()
			if err != nil {
				return portConfig{}, err
			}

			hostPort = port
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
	}, nil
}

// getRandomPort Returns an open, ephemeral port for mapping a container.
func getRandomPort() (string, error) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(handler)
	urlParts, err := url.ParseRequestURI(server.URL)

	server.Close()

	if err != nil {
		return "", err
	}

	return urlParts.Port(), nil
}
