package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ChrisWiegman/kana/internal/helpers"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Description string
	Command     []string
	Docker      bool
	Output      string
}

func Setup() {
	testDirectory := "./kana"

	dirExists, err := helpers.PathExists(testDirectory)
	if err != nil {
		panic(err)
	}

	if dirExists {
		err = os.RemoveAll(testDirectory)
		if err != nil {
			panic(err)
		}

		err = os.Mkdir(testDirectory, 0755) //nolint: mnd
		if err != nil {
			panic(err)
		}
	}

	err = os.Chdir(testDirectory)
	if err != nil {
		panic(err)
	}
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

	if docker {
		commands := [][]string{
			{"kill", "$(docker ps -q)"},
			{"container", "prune", "-f"},
			{"network", "prune", "-f"},
			{"volume", "prune", "-f"},
			{"system", "prune", "-a", "-f"}}

		for _, command := range commands {
			cmd := exec.Command("docker", command...)
			err = cmd.Run()
			if err != nil {
				panic(err)
			}
		}
	}
}

func RunCommandTest(testCases []Test, t *testing.T) {
	for _, test := range testCases {
		t.Run(test.Description, func(t *testing.T) {
			cmd := exec.Command("../../build/kana", test.Command...) //nolint: gosec
			var out, stdErr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stdErr
			err := cmd.Run()
			if err != nil {
				t.Fatalf("Unexpected error: %v", stdErr.String())
			}

			if test.Output == "" {
				snaps.MatchSnapshot(t, out.String())
			} else {
				assert.Contains(t, out.String(), test.Output)
			}

			Teardown(test.Docker)
		})
	}
}
