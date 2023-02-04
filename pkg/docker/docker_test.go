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

	c := new(DockerClient)

	c.mobyClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	execCommand = mocks.MockExecCommand
	mocks.MockedExitStatus = 1
	mocks.MockedStdout = "this is an error"
	defer func() { execCommand = exec.Command }()

	err = c.ensureDockerIsAvailable(consoleOutput)
	fmt.Println(err)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
