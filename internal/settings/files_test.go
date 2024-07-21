package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureKanaPlugin(t *testing.T) {
	siteDirectory := "."
	version := "1.0.0"
	siteName := "example.com"

	err := EnsureKanaPlugin(siteDirectory, version, siteName)
	if err != nil {
		t.Errorf("EnsureKanaPlugin returned an error: %v", err)
	}

	pluginPath := filepath.Join(siteDirectory, "wp-content", "mu-plugins", "kana-local-development.php")
	_, err = os.Stat(pluginPath)
	if err != nil {
		t.Errorf("Failed to create Kana plugin file: %v", err)
	}

	err = os.RemoveAll("./wp-content")
	if err != nil {
		t.Errorf("EnsureKanaPlugin returned an error: %v", err)
	}
}
func TestGetDefaultFilePermissions(t *testing.T) {
	dirPerms, filePerms := GetDefaultFilePermissions()

	expectedDirPerms := defaultDirPermissions
	expectedFilePerms := defaultFilePermissions

	if dirPerms != expectedDirPerms {
		t.Errorf("Incorrect default directory permissions. Got: %d, Expected: %d", dirPerms, expectedDirPerms)
	}

	if filePerms != expectedFilePerms {
		t.Errorf("Incorrect default file permissions. Got: %d, Expected: %d", filePerms, expectedFilePerms)
	}
}
func TestEnsureStaticConfigFiles(t *testing.T) {
	appDirectory := "."

	err := ensureStaticConfigFiles(appDirectory)
	if err != nil {
		t.Errorf("ensureStaticConfigFiles returned an error: %v", err)
	}

	for _, file := range configFiles {
		filePath := filepath.Join(appDirectory, file.LocalPath)
		destFile := filepath.Join(appDirectory, file.LocalPath, file.Name)

		_, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("Failed to create directory: %v", err)
		}

		_, err = os.Stat(destFile)
		if err != nil {
			t.Errorf("Failed to create file: %v", err)
		}
	}

	err = os.RemoveAll("./config")
	if err != nil {
		t.Errorf("EnsureKanaPlugin returned an error: %v", err)
	}
}
