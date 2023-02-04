package docker

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/ChrisWiegman/kana-cli/pkg/docker/mocks"

	"github.com/docker/docker/client"
)

func TestEnsureDockerIsAvailable(t *testing.T) {
	consoleOutput := new(console.Console)
	var err error

	dockerClient := new(DockerClient)

	dockerClient.moby, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	execCommand = mocks.MockExecCommand
	mocks.MockedExitStatus = 1
	mocks.MockedStdout = "this is an error"
	defer func() { execCommand = exec.Command }()

	err = ensureDockerIsAvailable(consoleOutput, dockerClient.moby)
	fmt.Println(err)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
