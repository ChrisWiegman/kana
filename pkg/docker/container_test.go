package docker

import (
	"testing"
)

func TestContainerRun(t *testing.T) {
	d, err := NewController()

	if err != nil {
		t.Error(err)
	}

	err = d.EnsureImage("alpine")
	if err != nil {
		t.Error(err)
	}

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
