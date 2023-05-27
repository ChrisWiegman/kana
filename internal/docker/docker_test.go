package docker

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"testing"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/docker/mocks"
	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestEnsureDockerIsAvailable(t *testing.T) {
	consoleOutput := new(console.Console)
	outputError := fmt.Errorf("we have an error")

	maxRetries = 1 // Only retry this once to save time.

	var tests = []struct {
		name           string
		goos           string
		dockerOutput   error
		expectedResult error
		exitStatus     int
	}{
		{
			"Test docker is running and no errors on list function.",
			"any",
			nil,
			nil,
			0},
		{
			"Test error on docker list function on Linux",
			"linux",
			outputError,
			outputError,
			0},
		{
			"Test error on first list function and unable to start Docker on Mac",
			"darwin",
			outputError,
			fmt.Errorf("error: unable to start Docker for Mac"),
			1},
	}

	for _, test := range tests {
		if test.goos == "darwin" && runtime.GOOS != "darwin" {
			continue
		}

		if test.goos == "linux" && runtime.GOOS != "linux" {
			continue
		}

		moby := new(mocks.APIClient)

		if test.exitStatus == 0 {
			moby.On("ContainerList", context.Background(), types.ContainerListOptions{}).Return([]types.Container{}, test.dockerOutput).Once()
		} else {
			moby.On("ContainerList", context.Background(), types.ContainerListOptions{}).Return([]types.Container{}, fmt.Errorf(""))
		}

		execCommand = mocks.MockExecCommand
		mocks.MockedExitStatus = test.exitStatus

		err := ensureDockerIsAvailable(consoleOutput, moby)
		assert.Equal(t, test.expectedResult, err, test.name)

		execCommand = exec.Command
	}
}
