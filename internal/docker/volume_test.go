package docker

import (
	"testing"
)

func TestSingleCreate(t *testing.T) {

	d, err := NewController()

	if err != nil {
		t.Error(err)
	}

	created, _, _ := d.EnsureVolume("myvolume")
	if created != true {
		t.Errorf("Should have created the volume the first time")
	}

	created, _, _ = d.EnsureVolume("myvolume")
	if created != false {
		t.Errorf("Should not have created the volume the second time")
	}

	removed, _ := d.RemoveVolume("myvolume")
	if removed != true {
		t.Errorf("Should have removed the volume")
	}
}

func TestEnsureVolume(t *testing.T) {

	d, err := NewController()

	if err != nil {
		t.Error(err)
	}

	_, volume, err := d.EnsureVolume("myvolume")

	if err != nil {
		t.Error(err)
	}

	if volume.Name != "myvolume" {
		t.Errorf("Expected volume name to be %s; got %s\n", "myvolume", volume.Name)
		t.FailNow()
	}

	removed, err := d.RemoveVolume("myvolume")

	if err != nil {
		t.Error(err)
	}

	if removed != true {
		t.Errorf("Volume should have been removed but wasn't")
	}

}
