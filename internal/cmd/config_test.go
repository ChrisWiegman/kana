package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/ChrisWiegman/kana/tests"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestMain(m *testing.M) {
	code := m.Run()
	tests.Teardown()
	os.Exit(code)
}

func TestConfig(t *testing.T) {
	kanaTests := []tests.Test{
		{
			Description: "Test the default config command",
			Command:     []string{"config"}},
		{
			Description: "Retrieve the PHP value from the config command",
			Command:     []string{"config", "php"}},
	}

	for _, test := range kanaTests {
		t.Run(test.Description, func(t *testing.T) {
			cmd := exec.Command("../../build/kana", test.Command...) //nolint: gosec
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			snaps.MatchSnapshot(t, out.String())
		})
	}
}
