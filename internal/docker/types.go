package docker

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	containerTypes "github.com/docker/docker/api/types/container"
	networkTypes "github.com/docker/docker/api/types/network"
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

// ContainerAPIClient defines API client methods for the containers.
type ContainerAPIClient interface {
	ContainerCreate(
		ctx context.Context,
		config *containerTypes.Config,
		hostConfig *containerTypes.HostConfig,
		networkingConfig *networkTypes.NetworkingConfig,
		platform *specs.Platform,
		containerName string) (containerTypes.CreateResponse, error)
	ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
	ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error)
	ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error)
	ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error)
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error
	ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error
	ContainerStop(ctx context.Context, name string, options containerTypes.StopOptions) error
	ContainerWait(
		ctx context.Context,
		container string,
		condition containerTypes.WaitCondition) (<-chan containerTypes.WaitResponse, <-chan error)
}

// ImageAPIClient defines API client methods for the images.
type ImageAPIClient interface {
	ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, image string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error)
	ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
}

// NetworkAPIClient defines API client methods for the networks.
type NetworkAPIClient interface {
	NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error)
	NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error)
	NetworkRemove(ctx context.Context, network string) error
}

// ViperClient defines a mock Viper client for testing.
type ViperClient interface {
	SetConfigName(in string)
	SetConfigType(in string)
	AddConfigPath(in string)
	ReadInConfig() error
	SafeWriteConfig() error
	GetTime(key string) time.Time
	Set(key string, value interface{})
	WriteConfig() error
}
