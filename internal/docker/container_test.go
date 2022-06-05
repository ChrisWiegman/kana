package docker

import (
	"testing"

	"github.com/ChrisWiegman/kana/internal/config"
)

func TestContainerRun(t *testing.T) {

	kanaConfig, err := config.GetKanaConfig()

	c, err := NewController(kanaConfig)

	if err != nil {
		t.Error(err)
	}

	err = c.EnsureImage("alpine")
	if err != nil {
		t.Error(err)
	}

	config := ContainerConfig{
		Image:   "alpine",
		Command: []string{"echo", "hello world"},
	}

	statusCode, body, err := c.ContainerRunAndClean(config)

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

	_, err = c.RemoveImage("alpine")
	if err != nil {
		t.Error(err)
	}
}
