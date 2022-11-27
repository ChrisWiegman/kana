package config

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed source/dynamic.toml
var DYNAMIC_TOML string

//go:embed source/traefik.toml
var TRAEFIK_TOML string

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
		LocalPath:   "config/traefik",
		Permissions: 0644,
		Overwrite:   true,
	},
	{
		Name:        "traefik.toml",
		Template:    TRAEFIK_TOML,
		LocalPath:   "config/traefik",
		Permissions: 0644,
		Overwrite:   true,
	},
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

	return os.WriteFile(destFile, finalTemplate, file.Permissions)
}
