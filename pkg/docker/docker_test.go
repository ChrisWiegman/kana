package docker

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"testing"

	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/ChrisWiegman/kana-cli/pkg/docker/mocks"
	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestEnsureDockerIsAvailable(t *testing.T) {
	consoleOutput := new(console.Console)

	var tests = []struct {
		goos       string
		output     error
		expected   error
		exitStatus int
		stdOut     string
	}{
		{"any", nil, nil, 0, ""},
		{"any", fmt.Errorf("we have an error"), fmt.Errorf("we have an error"), 0, ""},
		{"darwin", nil, fmt.Errorf("we have an error"), 1, "we have an error"},
	}

	for _, test := range tests {
		if test.goos == "darwin" && runtime.GOOS != "darwin" {
			continue
		}

		moby := new(mocks.APIClient)
		moby.On("ContainerList", context.Background(), types.ContainerListOptions{}).Return([]types.Container{}, test.output)

		execCommand = mocks.MockExecCommand
		mocks.MockedExitStatus = test.exitStatus
		mocks.MockedStdout = test.stdOut
		execCommand = exec.Command

		err := ensureDockerIsAvailable(consoleOutput, moby)
		assert.Equal(t, test.expected, err)
	}
}
