package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
)

type VolumeMount struct {
	HostPath string
	Volume   *types.Volume
}

func (c *Controller) FindVolume(name string) (volume *types.Volume, err error) {
	volumes, err := c.cli.VolumeList(context.Background(), filters.NewArgs())

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

func (c *Controller) EnsureVolume(name string) (created bool, volume *types.Volume, err error) {
	volume, err = c.FindVolume(name)

	if err != nil {
		return false, nil, err
	}

	if volume != nil {
		return false, volume, nil
	}

	vol, err := c.cli.VolumeCreate(context.Background(), volumetypes.VolumeCreateBody{
		Driver: "local",
		//		DriverOpts: map[string]string{},
		//		Labels:     map[string]string{},
		Name: name,
	})

	return true, &vol, err
}

func (c *Controller) RemoveVolume(name string) (removed bool, err error) {
	vol, err := c.FindVolume(name)

	if err != nil {
		return false, err
	}

	if vol == nil {
		return false, nil
	}

	err = c.cli.VolumeRemove(context.Background(), name, true)

	if err != nil {
		return false, err
	}

	return true, nil
}
