package mocks

import (
	io "io"

	"github.com/moby/moby/pkg/jsonmessage"
)

var MockedDisplayJSONMessagesStreamReturn error

func MockDisplayJSONMessagesStream(in io.Reader, out io.Writer, terminalFd uintptr, isTerminal bool, auxCallback func(jsonmessage.JSONMessage)) error { //nolint:lll
	return MockedDisplayJSONMessagesStreamReturn
}
