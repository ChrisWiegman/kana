package helpers

import (
	"archive/zip"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

// CopyFile copies a file from source to destination.
func CopyFile(sourceFile, destinationFile string) error {
	sourceFileStat, err := os.Stat(sourceFile)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", sourceFile)
	}

	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// DownloadFile downloads a file from a given URL and saves it to the destination path.
func DownloadFile(downloadURL, destinationPath string) (string, error) {
	// Build fileName from fullPath
	fileURL, err := url.Parse(downloadURL)
	if err != nil {
		return "", err
	}

	fileName := filepath.Base(fileURL.Path)

	// Create blank file
	file, err := os.Create(filepath.Join(destinationPath, fileName))
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, downloadURL, http.NoBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(file, resp.Body)
	defer file.Close()

	return fileName, err
}

// UnZipFile extracts a zip file to a given destination path.
func UnZipFile(sourceFile, destinationPath string) error {
	archive, err := zip.OpenReader(sourceFile)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(destinationPath, filepath.Clean(f.Name))

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, f.Mode())
			if err != nil {
				return err
			}

			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		const maxBufferSize = 1 << 20

		for {
			_, err := io.CopyN(dstFile, fileInArchive, maxBufferSize)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}
