package cmd

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestConfig(t *testing.T) {
	t.Run("Test the default config command", func(t *testing.T) {
		cmd := exec.Command("../../build/kana", "config")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		snaps.MatchSnapshot(t, out.String())
	})

	t.Run("Retrieve the PHP value from the config command", func(t *testing.T) {
		cmd := exec.Command("../../build/kana", "config", "php")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		snaps.MatchSnapshot(t, out.String())
	})
}
