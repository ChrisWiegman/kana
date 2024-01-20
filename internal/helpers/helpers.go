package helpers

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// IsValidString Checks a given string against an array of valid values and returns true/false as appropriate.
func IsValidString(stringToCheck string, validStrings []string) bool {
	for _, validString := range validStrings {
		if validString == stringToCheck {
			return true
		}
	}

	return false
}

// SanitizeSiteName Returns the site name, properly sanitized for use.
func SanitizeSiteName(rawSiteName string) string {
	siteName := strings.TrimSpace(rawSiteName)
	siteName = strings.ToLower(siteName)
	siteName = strings.ReplaceAll(siteName, " ", "-")
	siteName = strings.ReplaceAll(siteName, "_", "-")
	return strings.ToValidUTF8(siteName, "")
}

// PathExists returns true if the given path exists or false if it doesn't.
func PathExists(filePath string) (bool, error) {
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
func ReadLine(r *bufio.Reader) (string, error) {
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

// IsEmpty returns a bool to indicate if the provided path is empty.
func IsEmpty(path string) (bool, error) {
	osFile, err := os.Open(path)
	if err != nil {
		return false, err
	}

	defer osFile.Close()

	_, err = osFile.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}
