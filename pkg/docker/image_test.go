package docker

import "testing"

func TestEnsureImage(t *testing.T) {
	c, err := NewController()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = c.EnsureImage("alpine")

	if err != nil {
		t.Error(err)
	}
}
