package docker

import (
	"testing"
)

func TestContainerRun(t *testing.T) {
	c, err := NewController()

	if err != nil {
		t.Error(err)
	}

	statusCode, body, err := c.ContainerRunAndClean("alpine", []string{"echo", "hello world"}, []VolumeMount{})

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
}
