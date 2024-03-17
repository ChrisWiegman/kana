package cmd

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestRoot(t *testing.T) {
	t.Run("run the kana root command without further input", func(t *testing.T) {
		cmd := exec.Command("../../build/kana")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		snaps.MatchSnapshot(t, out.String())
	})
}
