package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// APIClient is an interface that clients that talk with a docker server must implement.
type APIClient interface {
	ContainerAPIClient
	ImageAPIClient
	NetworkAPIClient
}

// Ensure that Client always implements APIClient.
var _ APIClient = &client.Client{}

// ContainerAPIClient defines API client methods for the containers
type ContainerAPIClient interface {
	ContainerCreate(
		ctx context.Context,
		config *containertypes.Config,
		hostConfig *containertypes.HostConfig,
		networkingConfig *networktypes.NetworkingConfig,
		platform *specs.Platform,
		containerName string) (containertypes.CreateResponse, error)
	ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
	ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error)
	ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error)
	ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error)
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error
	ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error
	ContainerStop(ctx context.Context, name string, options containertypes.StopOptions) error
	ContainerWait(
		ctx context.Context,
		container string,
		condition containertypes.WaitCondition) (<-chan containertypes.WaitResponse, <-chan error)
}

// ImageAPIClient defines API client methods for the images
type ImageAPIClient interface {
	ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, image string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error)
}

// NetworkAPIClient defines API client methods for the networks
type NetworkAPIClient interface {
	NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error)
	NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error)
	NetworkRemove(ctx context.Context, network string) error
}
