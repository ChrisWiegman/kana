package docker

import (
	"testing"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/docker/mocks"

	"github.com/knadh/koanf/v2"
)

func TestContainerRun(t *testing.T) {
	consoleOutput := new(console.Console)

	d, err := New(consoleOutput, "")
	if err != nil {
		t.Error(err)
	}

	d.imageUpdateData = koanf.New(".")

	err = d.EnsureImage("alpine", 1, consoleOutput)
	if err != nil {
		t.Error(err)
	}

	displayJSONMessagesStream = mocks.MockDisplayJSONMessagesStream
	mocks.MockedDisplayJSONMessagesStreamReturn = nil //nolint:gocritic

	config := ContainerConfig{
		Image:   "alpine",
		Command: []string{"echo", "hello world"},
	}

	statusCode, body, err := d.ContainerRunAndClean(&config, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if body != "hello world\r\n" {
		t.Errorf("Expected 'hello world'; received %q\n", body)
	}

	if statusCode != 0 {
		t.Errorf("Expect status to be 0; received %q\n", statusCode)
	}

	_, err = d.removeImage("alpine")
	if err != nil {
		t.Error(err)
	}
}
