package cmd //nolint: dupl

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestList(t *testing.T) {
	t.Run("Test the default list command", func(t *testing.T) {
		cmd := exec.Command("../../build/kana", "list")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		snaps.MatchSnapshot(t, out.String())
	})

	t.Run("Test the list command with json output", func(t *testing.T) {
		cmd := exec.Command("../../build/kana", "list", "--output-json")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		snaps.MatchSnapshot(t, out.String())
	})
}
