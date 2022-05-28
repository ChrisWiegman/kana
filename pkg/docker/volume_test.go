package docker

import (
	"testing"
)

func TestSingleCreate(t *testing.T) {
	c, err := NewController()

	if err != nil {
		t.Error(err)
	}

	created, _, _ := c.EnsureVolume("myvolume")
	if created != true {
		t.Errorf("Should have created the volume the first time")
	}

	created, _, _ = c.EnsureVolume("myvolume")
	if created != false {
		t.Errorf("Should not have created the volume the second time")
	}

	removed, _ := c.RemoveVolume("myvolume")
	if removed != true {
		t.Errorf("Should have removed the volume")
	}
}

func TestEnsureVolume(t *testing.T) {
	c, err := NewController()

	if err != nil {
		t.Error(err)
	}

	_, volume, err := c.EnsureVolume("myvolume")

	if err != nil {
		t.Error(err)
	}

	if volume.Name != "myvolume" {
		t.Errorf("Expected volume name to be %s; got %s\n", "myvolume", volume.Name)
		t.FailNow()
	}

	removed, err := c.RemoveVolume("myvolume")

	if err != nil {
		t.Error(err)
	}

	if removed != true {
		t.Errorf("Volume should have been removed but wasn't")
	}

}

func TestPersistentVolume(t *testing.T) {
	c, err := NewController()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	created, volume, err := c.EnsureVolume("persistentvolume")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if created != true {
		t.Errorf("Should have created a volume at the start")
	}

	mounts := []VolumeMount{
		{
			HostPath: "/volume",
			Volume:   volume,
		},
	}

	_, body1, _ := c.ContainerRunAndClean("testimage", []string{}, mounts)

	// Second run

	statusCode, body2, err := c.ContainerRunAndClean("testimage", []string{}, mounts)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if statusCode != 0 {
		t.Error("Second run should not have created a file")
	}

	if body1 != body2 {
		t.Errorf("%s\nShould have been equal to:\n%s\n", body1, body2)
	}

	c.RemoveVolume("persistentvolume")
}
