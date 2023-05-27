package docker

import (
	"testing"
	"time"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/internal/docker/mocks"
	"github.com/stretchr/testify/mock"
)

func TestContainerRun(t *testing.T) {
	consoleOutput := new(console.Console)

	d, err := NewDockerClient(consoleOutput, "")
	if err != nil {
		t.Error(err)
	}

	viper := new(mocks.ViperClient)
	viper.On("ReadInConfig").Return(nil)
	viper.On("GetTime", mock.Anything).Return(time.Now())
	viper.On("Set", mock.Anything, mock.Anything).Return()
	viper.On("WriteConfig").Return(nil)

	d.imageUpdateData = viper

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

	statusCode, body, err := d.ContainerRunAndClean(&config)
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

	_, err = d.RemoveImage("alpine")
	if err != nil {
		t.Error(err)
	}
}
