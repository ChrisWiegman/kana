package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/docker/docker/client"
)

var mockedExitStatus = 0
var mockedStdout string

func TestEnsureDockerIsAvailable(t *testing.T) {
	consoleOutput := new(console.Console)
	var err error

	c := new(DockerClient)

	c.client, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	execCommand = mockExecCommand
	mockedExitStatus = 1
	mockedStdout = "this is an error"
	defer func() { execCommand = exec.Command }()

	err = c.ensureDockerIsAvailable(consoleOutput)
	fmt.Println(err)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func mockExecCommand(command string, args ...string) *exec.Cmd {
	fmt.Println("mocking")
	cs := []string{"-test.run=TestExecCommandHelper", "--", command}
	cs = append(cs, args...)

	cmd := exec.Command(os.Args[0], cs...) //nolint:gosec
	es := strconv.Itoa(mockedExitStatus)

	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1",
		"STDOUT=" + mockedStdout,
		"EXIT_STATUS=" + es}

	return cmd
}
