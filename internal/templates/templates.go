package templates

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type File struct {
	Name, Template, LocalPath string
	Permissions               os.FileMode
	Replacements              []Replacement
	Overwrite                 bool
}

// Replacement Replacement struct includes templating search-replace information.
type Replacement struct {
	Search, Replace string
	Count           int
}

var configFiles = []File{
	{
		Name:        "dynamic.toml",
		Template:    DYNAMIC_TOML,
		LocalPath:   "conf/traefik",
		Permissions: 0644,
		Overwrite:   true,
	},
	{
		Name:        "traefik.toml",
		Template:    TRAEFIK_TOML,
		LocalPath:   "conf/traefik",
		Permissions: 0644,
		Overwrite:   true,
	},
}

// GetConfigRoot Return the path for the global config.
func GetConfigRoot() (string, error) {

	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".kana"), nil

}

// WriteConfigFiles Write config files to an install's .wpengine-conf directory
func WriteConfigFiles() error {

	// Get global config root.
	configRoot, err := GetConfigRoot()
	if err != nil {
		return err
	}

	return writeFileArrayToDisk(configFiles, configRoot)

}

func writeFileArrayToDisk(files []File, installPath string) error {

	for _, file := range files {

		// Don't overwrite the file if Overwrite is false and the file exists.
		destFile := filepath.Join(installPath, file.LocalPath, file.Name)

		if !file.Overwrite {
			_, err := os.Stat(destFile)
			if !os.IsNotExist(err) {
				continue
			}
		}

		if err := writeFileFromTemplate(installPath, file); err != nil {
			return err
		}
	}

	return nil

}

func writeFileFromTemplate(installPath string, file File) error {

	filePath := filepath.Join(installPath, file.LocalPath)
	destFile := filepath.Join(installPath, file.LocalPath, file.Name)

	if err := os.MkdirAll(filePath, 0750); err != nil {
		return err
	}

	finalTemplate := []byte(file.Template)

	if len(file.Replacements) > 0 {

		for _, replacement := range file.Replacements {
			finalTemplate = bytes.Replace(finalTemplate, []byte(replacement.Search), []byte(replacement.Replace), replacement.Count)
		}
	}

	return ioutil.WriteFile(destFile, finalTemplate, file.Permissions)

}
