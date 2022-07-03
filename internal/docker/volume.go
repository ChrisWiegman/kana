package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
)

func (d *DockerClient) FindVolume(name string) (volume *types.Volume, err error) {
	volumes, err := d.client.VolumeList(context.Background(), filters.NewArgs())

	if err != nil {
		return nil, err
	}

	for _, v := range volumes.Volumes {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, nil
}

func (d *DockerClient) EnsureVolume(name string) (created bool, volume *types.Volume, err error) {

	volume, err = d.FindVolume(name)

	if err != nil {
		return false, nil, err
	}

	if volume != nil {
		return false, volume, nil
	}

	vol, err := d.client.VolumeCreate(context.Background(), volumetypes.VolumeCreateBody{
		Driver: "local",
		//		DriverOpts: map[string]string{},
		//		Labels:     map[string]string{},
		Name: name,
	})

	return true, &vol, err
}

func (d *DockerClient) RemoveVolume(name string) (removed bool, err error) {

	vol, err := d.FindVolume(name)

	if err != nil {
		return false, err
	}

	if vol == nil {
		return false, nil
	}

	err = d.client.VolumeRemove(context.Background(), name, true)

	if err != nil {
		return false, err
	}

	return true, nil
}
