package docker

import (
	"testing"

	"github.com/ChrisWiegman/kana-cli/pkg/console"
)

func TestEnsureImage(t *testing.T) {
	consoleOutput := new(console.Console)

	d, err := NewDockerClient(consoleOutput)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = d.EnsureImage("alpine", consoleOutput)

	if err != nil {
		t.Error(err)
	}
}

func TestRemoveImage(t *testing.T) {
	consoleOutput := new(console.Console)

	d, err := NewDockerClient(consoleOutput)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = d.EnsureImage("alpine", consoleOutput)

	if err != nil {
		t.Error(err)
	}

	removed, err := d.RemoveImage("alpine")
	if err != nil {
		t.Error(err)
	}

	if removed != true {
		t.Errorf("Image should have been removed but wasn't")
	}

	removed, err = d.RemoveImage("alpine")
	if err != nil {
		t.Error(err)
	}

	if removed == true {
		t.Errorf("Image should not have been removed but was")
	}
}
