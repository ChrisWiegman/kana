package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestGetStaticDirectories(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testDir := filepath.Join(cwd, "kana-test")

	// Create a temporary directory for testing
	err = os.Mkdir(testDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	// Set the current working directory to the temporary directory
	err = os.Chdir(testDir)
	if err != nil {
		t.Fatal(err)
	}

	// Call the function being tested
	directories, err := getStaticDirectories()
	if err != nil {
		t.Fatal(err)
	}

	// Assert the expected values
	expectedWorking := testDir

	if directories.Working != expectedWorking {
		t.Errorf("Expected Working directory to be %s, but got %s", expectedWorking, directories.Working)
	}

	homeDir, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}
	expectedApp := filepath.Join(homeDir, configFolderName)
	if directories.App != expectedApp {
		t.Errorf("Expected App directory to be %s, but got %s", expectedApp, directories.App)
	}
}
