package site

import (
	"bufio"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/ChrisWiegman/kana/internal/helpers"
)

const DefaultType = "site"

func (s *Site) DetectType() (string, error) {
	var err error
	var isSite bool

	isSite, err = helpers.PathExists(path.Join(s.Settings.WorkingDirectory, "wp-includes", "version.php"))
	if err != nil {
		return "", err
	}

	if isSite {
		return DefaultType, err
	}

	items, _ := os.ReadDir(s.Settings.WorkingDirectory)

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		if item.Name() == "style.css" || filepath.Ext(item.Name()) == ".php" {
			var f *os.File
			var line string

			f, err = os.Open(path.Join(s.Settings.WorkingDirectory, item.Name()))
			if err != nil {
				return "", err
			}

			reader := bufio.NewReader(f)
			line, err = helpers.ReadLine(reader)

			for err == nil {
				exp := regexp.MustCompile(`(Plugin|Theme) Name: .*`)

				for _, match := range exp.FindAllStringSubmatch(line, -1) {
					if match[1] == "Theme" {
						return "theme", err //nolint
					} else {
						return "plugin", err //nolint
					}
				}
				line, err = helpers.ReadLine(reader)
			}
		}
	}

	// We don't care if it is an empty folder.
	if err == io.EOF {
		err = nil
	}

	return DefaultType, err
}
