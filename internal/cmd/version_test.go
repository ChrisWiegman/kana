package cmd

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestVersion(t *testing.T) {
	t.Run("Test the version command for appropriate output", func(t *testing.T) {
		cmd := exec.Command("../../build/kana", "version")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		snaps.MatchSnapshot(t, out.String())
	})
}
