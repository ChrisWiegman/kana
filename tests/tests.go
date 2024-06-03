package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/mitchellh/go-homedir"
)

type Test struct {
	Description string
	Command     []string
	Docker      bool
}

func Teardown(docker bool) {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	appDirectory := filepath.Join(home, ".config", "kana")

	err = os.RemoveAll(appDirectory)
	if err != nil {
		panic(err)
	}
}

func RunSnapshotTest(testCases []Test, t *testing.T) {
	for _, test := range testCases {
		t.Run(test.Description, func(t *testing.T) {
			cmd := exec.Command("../../build/kana", test.Command...) //nolint: gosec
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			snaps.MatchSnapshot(t, out.String())

			Teardown(test.Docker)
		})
	}
}
