package site

import (
	"fmt"
	"io"
	"os"
)

// arrayContains Searches an array of strings for a given string and returns true/false as appropriate.
func arrayContains(array []string, name string) bool {
	for _, value := range array {
		if value == name {
			return true
		}
	}

	return false
}

// copyFile Copies a file on the user's host from one place to another.
func copyFile(src, dest string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("please enter a valid sql file")
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}

	defer func() {
		if err = source.Close(); err != nil {
			panic(err)
		}
	}()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer func() {
		if err = destination.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(destination, source)
	return err
}
