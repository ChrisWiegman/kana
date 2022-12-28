package docker

import (
	"testing"
)

func TestNetworkCreate(t *testing.T) {
	d, err := NewController()

	if err != nil {
		t.Error(err)
	}

	created, _, _ := d.EnsureNetwork("mynetwork")
	if created != true {
		t.Errorf("Should have created the network the first time")
	}

	created, _, _ = d.EnsureNetwork("mynetwork")
	if created != false {
		t.Errorf("Should not have created the network the second time")
	}

	removed, _ := d.RemoveNetwork("mynetwork")
	if removed != true {
		t.Errorf("Should have removed the network")
	}
}

func TestEnsureNetwork(t *testing.T) {
	d, err := NewController()

	if err != nil {
		t.Error(err)
	}

	_, network, err := d.EnsureNetwork("mynetwork")

	if err != nil {
		t.Error(err)
	}

	if network.Name != "mynetwork" {
		t.Errorf("Expected network name to be %s; got %s\n", "mynetwork", network.Name)
		t.FailNow()
	}

	removed, err := d.RemoveNetwork("mynetwork")

	if err != nil {
		t.Error(err)
	}

	if removed != true {
		t.Errorf("Network should have been removed but wasn't")
	}
}
