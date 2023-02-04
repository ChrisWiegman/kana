package mocks

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var MockedExitStatus = 0
var MockedStdout string

func MockExecCommand(command string, args ...string) *exec.Cmd {
	fmt.Println("mocking")
	cs := []string{"-test.run=TestExecCommandHelper", "--", command}
	cs = append(cs, args...)

	cmd := exec.Command(os.Args[0], cs...) //nolint:gosec
	es := strconv.Itoa(MockedExitStatus)

	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1",
		"STDOUT=" + MockedStdout,
		"EXIT_STATUS=" + es}

	return cmd
}
