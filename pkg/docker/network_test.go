package docker

import (
	"testing"
)

func TestNetworkCreate(t *testing.T) {
	c, err := NewController()

	if err != nil {
		t.Error(err)
	}

	created, _, _ := c.EnsureNetwork("mynetwork")
	if created != true {
		t.Errorf("Should have created the network the first time")
	}

	created, _, _ = c.EnsureNetwork("mynetwork")
	if created != false {
		t.Errorf("Should not have created the network the second time")
	}

	removed, _ := c.RemoveNetwork("mynetwork")
	if removed != true {
		t.Errorf("Should have removed the network")
	}
}

func TestEnsureNetwork(t *testing.T) {
	c, err := NewController()

	if err != nil {
		t.Error(err)
	}

	_, network, err := c.EnsureNetwork("mynetwork")

	if err != nil {
		t.Error(err)
	}

	if network.Name != "mynetwork" {
		t.Errorf("Expected network name to be %s; got %s\n", "mynetwork", network.Name)
		t.FailNow()
	}

	removed, err := c.RemoveNetwork("mynetwork")

	if err != nil {
		t.Error(err)
	}

	if removed != true {
		t.Errorf("Network should have been removed but wasn't")
	}

}
