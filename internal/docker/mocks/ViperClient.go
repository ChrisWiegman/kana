// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// ViperClient is an autogenerated mock type for the ViperClient type
type ViperClient struct {
	mock.Mock
}

// AddConfigPath provides a mock function with given fields: in
func (_m *ViperClient) AddConfigPath(in string) {
	_m.Called(in)
}

// GetTime provides a mock function with given fields: key
func (_m *ViperClient) GetTime(key string) time.Time {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for GetTime")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func(string) time.Time); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// ReadInConfig provides a mock function with given fields:
func (_m *ViperClient) ReadInConfig() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadInConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SafeWriteConfig provides a mock function with given fields:
func (_m *ViperClient) SafeWriteConfig() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SafeWriteConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Set provides a mock function with given fields: key, value
func (_m *ViperClient) Set(key string, value interface{}) {
	_m.Called(key, value)
}

// SetConfigName provides a mock function with given fields: in
func (_m *ViperClient) SetConfigName(in string) {
	_m.Called(in)
}

// SetConfigType provides a mock function with given fields: in
func (_m *ViperClient) SetConfigType(in string) {
	_m.Called(in)
}

// WriteConfig provides a mock function with given fields:
func (_m *ViperClient) WriteConfig() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for WriteConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewViperClient creates a new instance of ViperClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewViperClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ViperClient {
	mock := &ViperClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
