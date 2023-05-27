// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	context "context"

	container "github.com/docker/docker/api/types/container"

	io "io"

	mock "github.com/stretchr/testify/mock"

	network "github.com/docker/docker/api/types/network"

	types "github.com/docker/docker/api/types"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// ContainerAPIClient is an autogenerated mock type for the ContainerAPIClient type
type ContainerAPIClient struct {
	mock.Mock
}

// ContainerCreate provides a mock function with given fields: ctx, config, hostConfig, networkingConfig, platform, containerName
func (_m *ContainerAPIClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.CreateResponse, error) {
	ret := _m.Called(ctx, config, hostConfig, networkingConfig, platform, containerName)

	var r0 container.CreateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, *v1.Platform, string) (container.CreateResponse, error)); ok {
		return rf(ctx, config, hostConfig, networkingConfig, platform, containerName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, *v1.Platform, string) container.CreateResponse); ok {
		r0 = rf(ctx, config, hostConfig, networkingConfig, platform, containerName)
	} else {
		r0 = ret.Get(0).(container.CreateResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, *v1.Platform, string) error); ok {
		r1 = rf(ctx, config, hostConfig, networkingConfig, platform, containerName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerExecAttach provides a mock function with given fields: ctx, execID, config
func (_m *ContainerAPIClient) ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
	ret := _m.Called(ctx, execID, config)

	var r0 types.HijackedResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ExecStartCheck) (types.HijackedResponse, error)); ok {
		return rf(ctx, execID, config)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ExecStartCheck) types.HijackedResponse); ok {
		r0 = rf(ctx, execID, config)
	} else {
		r0 = ret.Get(0).(types.HijackedResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.ExecStartCheck) error); ok {
		r1 = rf(ctx, execID, config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerExecCreate provides a mock function with given fields: ctx, _a1, config
func (_m *ContainerAPIClient) ContainerExecCreate(ctx context.Context, _a1 string, config types.ExecConfig) (types.IDResponse, error) {
	ret := _m.Called(ctx, _a1, config)

	var r0 types.IDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ExecConfig) (types.IDResponse, error)); ok {
		return rf(ctx, _a1, config)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ExecConfig) types.IDResponse); ok {
		r0 = rf(ctx, _a1, config)
	} else {
		r0 = ret.Get(0).(types.IDResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.ExecConfig) error); ok {
		r1 = rf(ctx, _a1, config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerExecInspect provides a mock function with given fields: ctx, execID
func (_m *ContainerAPIClient) ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error) {
	ret := _m.Called(ctx, execID)

	var r0 types.ContainerExecInspect
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (types.ContainerExecInspect, error)); ok {
		return rf(ctx, execID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) types.ContainerExecInspect); ok {
		r0 = rf(ctx, execID)
	} else {
		r0 = ret.Get(0).(types.ContainerExecInspect)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, execID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerInspect provides a mock function with given fields: ctx, _a1
func (_m *ContainerAPIClient) ContainerInspect(ctx context.Context, _a1 string) (types.ContainerJSON, error) {
	ret := _m.Called(ctx, _a1)

	var r0 types.ContainerJSON
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (types.ContainerJSON, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) types.ContainerJSON); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(types.ContainerJSON)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerList provides a mock function with given fields: ctx, options
func (_m *ContainerAPIClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	ret := _m.Called(ctx, options)

	var r0 []types.Container
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.ContainerListOptions) ([]types.Container, error)); ok {
		return rf(ctx, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.ContainerListOptions) []types.Container); ok {
		r0 = rf(ctx, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Container)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.ContainerListOptions) error); ok {
		r1 = rf(ctx, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerLogs provides a mock function with given fields: ctx, _a1, options
func (_m *ContainerAPIClient) ContainerLogs(ctx context.Context, _a1 string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	ret := _m.Called(ctx, _a1, options)

	var r0 io.ReadCloser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ContainerLogsOptions) (io.ReadCloser, error)); ok {
		return rf(ctx, _a1, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ContainerLogsOptions) io.ReadCloser); ok {
		r0 = rf(ctx, _a1, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.ContainerLogsOptions) error); ok {
		r1 = rf(ctx, _a1, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerRemove provides a mock function with given fields: ctx, _a1, options
func (_m *ContainerAPIClient) ContainerRemove(ctx context.Context, _a1 string, options types.ContainerRemoveOptions) error {
	ret := _m.Called(ctx, _a1, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ContainerRemoveOptions) error); ok {
		r0 = rf(ctx, _a1, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerStart provides a mock function with given fields: ctx, _a1, options
func (_m *ContainerAPIClient) ContainerStart(ctx context.Context, _a1 string, options types.ContainerStartOptions) error {
	ret := _m.Called(ctx, _a1, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ContainerStartOptions) error); ok {
		r0 = rf(ctx, _a1, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerStop provides a mock function with given fields: ctx, name, options
func (_m *ContainerAPIClient) ContainerStop(ctx context.Context, name string, options container.StopOptions) error {
	ret := _m.Called(ctx, name, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, container.StopOptions) error); ok {
		r0 = rf(ctx, name, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerWait provides a mock function with given fields: ctx, _a1, condition
func (_m *ContainerAPIClient) ContainerWait(ctx context.Context, _a1 string, condition container.WaitCondition) (<-chan container.WaitResponse, <-chan error) {
	ret := _m.Called(ctx, _a1, condition)

	var r0 <-chan container.WaitResponse
	var r1 <-chan error
	if rf, ok := ret.Get(0).(func(context.Context, string, container.WaitCondition) (<-chan container.WaitResponse, <-chan error)); ok {
		return rf(ctx, _a1, condition)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, container.WaitCondition) <-chan container.WaitResponse); ok {
		r0 = rf(ctx, _a1, condition)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan container.WaitResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, container.WaitCondition) <-chan error); ok {
		r1 = rf(ctx, _a1, condition)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(<-chan error)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewContainerAPIClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewContainerAPIClient creates a new instance of ContainerAPIClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewContainerAPIClient(t mockConstructorTestingTNewContainerAPIClient) *ContainerAPIClient {
	mock := &ContainerAPIClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
