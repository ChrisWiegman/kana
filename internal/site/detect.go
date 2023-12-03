package site

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const defaultType = "site"

func (s *Site) DetectType() (string, error) {
	var err error
	var isSite bool

	isSite, err = pathExists(path.Join(s.Settings.WorkingDirectory, "wp-includes", "version.php"))
	if err != nil {
		return "", err
	}

	if isSite {
		return defaultType, err
	}

	items, _ := os.ReadDir(s.Settings.WorkingDirectory)

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		if item.Name() == "style.css" || filepath.Ext(item.Name()) == ".php" {
			var f *os.File
			var line string

			f, err = os.Open(item.Name())
			if err != nil {
				return "", err
			}

			reader := bufio.NewReader(f)
			line, err = readLine(reader)

			for err == nil {
				exp := regexp.MustCompile(`(Plugin|Theme) Name: .*`)

				for _, match := range exp.FindAllStringSubmatch(line, -1) {
					if match[1] == "Theme" {
						return "theme", err //nolint
					} else {
						return "plugin", err //nolint
					}
				}
				line, err = readLine(reader)
			}
		}
	}

	return defaultType, err
}

func pathExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)

	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// readLine returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix       = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
