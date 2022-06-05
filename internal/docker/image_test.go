package docker

import (
	"testing"

	"github.com/ChrisWiegman/kana/internal/config"
)

func TestEnsureImage(t *testing.T) {

	kanaConfig, err := config.GetKanaConfig()
	if err != nil {
		t.Error(err)
	}

	c, err := NewController(kanaConfig)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = c.EnsureImage("alpine")

	if err != nil {
		t.Error(err)
	}
}

func TestRemoveImage(t *testing.T) {

	kanaConfig, err := config.GetKanaConfig()
	if err != nil {
		t.Error(err)
	}

	c, err := NewController(kanaConfig)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = c.EnsureImage("alpine")

	if err != nil {
		t.Error(err)
	}

	removed, err := c.RemoveImage("alpine")
	if err != nil {
		t.Error(err)
	}

	if removed != true {
		t.Errorf("Image should have been removed but wasn't")
	}

	removed, err = c.RemoveImage("alpine")
	if err != nil {
		t.Error(err)
	}

	if removed == true {
		t.Errorf("Image should not have been removed but was")
	}
}
