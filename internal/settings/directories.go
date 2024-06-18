package settings

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func getStaticDirectories() (Directories, error) {
	directories := Directories{}

	cwd, err := os.Getwd()
	if err != nil {
		return directories, err
	}

	directories.Working = cwd

	home, err := homedir.Dir()
	if err != nil {
		return directories, err
	}

	directories.App = filepath.Join(home, configFolderName)

	return directories, nil
}
