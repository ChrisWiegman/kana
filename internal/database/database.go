package database

import (
	"fmt"
	"io"
	"os"
)

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
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
