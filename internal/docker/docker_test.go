package docker

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/ChrisWiegman/kana/internal/docker/mocks"

	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
)

func TestEnsureDockerIsAvailable(t *testing.T) {
	outputError := fmt.Errorf("Could not connect to Docker. Is Docker running?")

	var tests = []struct {
		name           string
		dockerOutput   error
		expectedResult error
		exitStatus     int
	}{
		{
			"Test docker is running and no errors on list function.",
			nil,
			nil,
			0},
		{
			"Test error on docker list function on Linux",
			outputError,
			outputError,
			0},
	}

	for _, test := range tests {
		apiClient := new(mocks.APIClient)

		if test.exitStatus == 0 {
			apiClient.On("ContainerList", context.Background(), container.ListOptions{}).Return([]container.Summary{}, test.dockerOutput).Once()
		} else {
			apiClient.On("ContainerList", context.Background(), container.ListOptions{}).Return([]container.Summary{}, fmt.Errorf(""))
		}

		execCommand = mocks.MockExecCommand
		mocks.MockedExitStatus = test.exitStatus

		err := ensureDockerIsAvailable(apiClient)
		assert.Equal(t, test.expectedResult, err, test.name)

		execCommand = exec.Command
	}
}
