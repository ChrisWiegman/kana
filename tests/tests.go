package tests

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type Test struct {
	Description string
	Command     []string
}

func Teardown() {
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
