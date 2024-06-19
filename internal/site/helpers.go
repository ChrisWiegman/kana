package site

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ChrisWiegman/kana/internal/docker"
)

// checkStatusCode returns true on 200 or false.
func checkStatusCode(checkURL string) (bool, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, checkURL, http.NoBody)
	if err != nil {
		return false, err
	}

	// Ignore SSL check as we're using our self-signed cert for development
	clientTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
	}

	client := &http.Client{
		Transport: clientTransport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusFound {
		return true, nil
	}

	return false, nil
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

// handleImageError Handles errors related to image detection and provides more helpful error messages.
func (s *Site) handleImageError(container *docker.ContainerConfig, err error) error {
	if strings.Contains(err.Error(), "manifest unknown") {
		switch container.Labels["kana.type"] {
		case "wordpress":
			return fmt.Errorf(
				"the PHP version in your configuration, %s, is invalid. See https://hub.docker.com/_/wordpress for a list of supported versions",
				s.settings.Get("PHP"))
		case "database":
			databaseURL := "https://hub.docker.com/_/mariadb"

			if s.settings.Get("Database") == "mysql" {
				databaseURL = "https://hub.docker.com/_/mysql"
			}

			return fmt.Errorf(
				"the database version in your configuration, %s, is invalid. See %s for a list of supported versions",
				s.settings.Get("DatabaseVersion"), databaseURL)
		}
	}

	return err
}
