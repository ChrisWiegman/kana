package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureKanaPlugin(t *testing.T) {
	siteDirectory := "kana-test"
	version := "1.0.0"
	siteName := "example.com"

	// Create a temporary directory for testing
	err := os.Mkdir(siteDirectory, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(siteDirectory)

	err = EnsureKanaPlugin(siteDirectory, version, siteName)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	pluginPath := filepath.Join(siteDirectory, "wp-content", "mu-plugins", "kana-local-development.php")
	_, err = os.Stat(pluginPath)
	if err != nil {
		t.Errorf("Plugin file not created: %v", err)
	}

	// Additional assertions can be added here to validate the contents of the plugin file.
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
	appDirectory := "kana-test"

	// Create a temporary directory for testing
	err := os.Mkdir(appDirectory, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(appDirectory)

	err = ensureStaticConfigFiles(appDirectory)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	for _, file := range configFiles {
		filePath := filepath.Join(appDirectory, file.LocalPath)
		destFile := filepath.Join(appDirectory, file.LocalPath, file.Name)

		_, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("Directory not created: %v", err)
		}

		_, err = os.Stat(destFile)
		if err != nil {
			t.Errorf("File not created: %v", err)
		}
	}
}
