package helpers

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidString(t *testing.T) {
	var testCases = []struct {
		name         string
		checkString  string
		validStrings []string
		shouldPass   bool
	}{
		{
			name:         "Ensure a valid string is valid",
			checkString:  "test",
			validStrings: []string{"test", "test2"},
			shouldPass:   true},
		{
			name:         "Ensure an invalid string is not valid",
			checkString:  "test",
			validStrings: []string{"test2", "test3"},
			shouldPass:   false},
	}

	for _, test := range testCases {
		result := IsValidString(test.checkString, test.validStrings)

		assert.Equal(t, test.shouldPass, result, test.name)
	}
}

func TestSanitizeSiteName(t *testing.T) {
	var testCases = []struct {
		name         string
		rawSiteName  string
		expectedName string
	}{
		{
			name:         "Ensure a simple name is sanitized",
			rawSiteName:  "Test Site",
			expectedName: "test-site"},
		{
			name:         "Ensure a complex name is sanitized",
			rawSiteName:  "Test_Site ",
			expectedName: "test-site"},
	}

	for _, test := range testCases {
		result := SanitizeSiteName(test.rawSiteName)

		assert.Equal(t, test.expectedName, result, test.name)
	}
}

func TestPathExists(t *testing.T) {
	testPath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	existingPath := testPath
	nonExistingPath := filepath.Join(testPath, "nonexisting")

	exists, err := PathExists(existingPath)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected existing path to return true")

	exists, err = PathExists(nonExistingPath)
	assert.NoError(t, err)
	assert.False(t, exists, "Expected non-existing path to return false")
}

func TestReadLine(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Write test data to the file
	data := "Line 1\nLine 2\nLine 3\n"
	_, err = file.WriteString(data)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	// Open the file for reading
	file, err = os.Open(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Create a buffered reader
	reader := bufio.NewReader(file)

	// Test reading each line
	expectedLines := []string{"Line 1", "Line 2", "Line 3"}
	for _, expectedLine := range expectedLines {
		line, err := ReadLine(reader)
		if err != nil {
			t.Fatal(err)
		}
		if line != expectedLine {
			t.Errorf("Expected line: %s, got: %s", expectedLine, line)
		}
	}
}
func TestIsEmpty(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Test with an empty directory
	emptyDir := filepath.Join(tempDir, "empty")

	err = os.Mkdir(emptyDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	isEmpty, err := IsEmpty(emptyDir)
	assert.NoError(t, err)
	assert.True(t, isEmpty, "Expected empty directory to return true")

	// Test with a non-empty directory
	nonEmptyDir := filepath.Join(tempDir, "nonempty")

	err = os.Mkdir(nonEmptyDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	nonEmptyFile := filepath.Join(nonEmptyDir, "file.txt")

	err = os.WriteFile(nonEmptyFile, []byte("test"), os.ModePerm) //nolint: gosec
	if err != nil {
		t.Fatal(err)
	}
	isEmpty, err = IsEmpty(nonEmptyDir)

	assert.NoError(t, err)
	assert.False(t, isEmpty, "Expected non-empty directory to return false")

	// Test with a non-existing directory
	nonExistingDir := filepath.Join(tempDir, "nonexisting")

	isEmpty, err = IsEmpty(nonExistingDir)

	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.False(t, isEmpty, "Expected non-existing directory to return false")
}

func TestCopyFile(t *testing.T) {
	testPath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	destinationFile := "file.txt"

	// Create a temporary source file for testing
	source, err := os.CreateTemp(testPath, "sourcefile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(source.Name())
	defer source.Close()

	// Write test data to the source file
	data := "Test data"
	_, err = source.WriteString(data)
	if err != nil {
		t.Fatal(err)
	}

	err = CopyFile(source.Name(), filepath.Join(testPath, destinationFile))
	if err != nil {
		t.Fatal(err)
	}

	// Read the destination file
	destination, err := os.Open(destinationFile)
	if err != nil {
		t.Fatal(err)
	}
	defer destination.Close()

	// Read the contents of the destination file
	destinationData, err := io.ReadAll(destination)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the contents of the destination file match the source file
	if string(destinationData) != data {
		t.Errorf("Expected destination file contents: %s, got: %s", data, string(destinationData))
	}

	// Delete the destination file
	err = os.Remove(destinationFile)
	if err != nil {
		t.Fatal(err)
	}
}
func TestDownloadFile(t *testing.T) {
	destinationPath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	fileName := "file.txt"
	fileContent := "Test data"

	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fileContent)
	}))
	defer mockServer.Close()

	// Build the download URL
	downloadURL, err := url.JoinPath(mockServer.URL, fileName)
	if err != nil {
		t.Fatal(err)
	}

	// Perform the download
	downloadedFileName, err := DownloadFile(downloadURL, destinationPath)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the downloaded file name matches the expected file name
	if downloadedFileName != fileName {
		t.Errorf("Expected downloaded file name: %s, got: %s", fileName, downloadedFileName)
	}

	// Read the downloaded file
	downloadedFilePath := filepath.Join(destinationPath, downloadedFileName)
	downloadedFileContent, err := os.ReadFile(downloadedFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the downloaded file content matches the expected file content
	if string(downloadedFileContent) != fileContent {
		t.Errorf("Expected downloaded file content: %s, got: %s", fileContent, string(downloadedFileContent))
	}

	// Clean up the downloaded file
	err = os.Remove(downloadedFilePath)
	if err != nil {
		t.Fatal(err)
	}
}
func TestUnZipFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary zip file for testing
	zipFile := filepath.Join(tempDir, "test.zip")
	err = createTempZipFile(zipFile)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(zipFile)

	// Extract the zip file
	err = UnZipFile(zipFile, tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the extracted files exist
	extractedFile1 := filepath.Join(tempDir, "file1.txt")
	exists, err := PathExists(extractedFile1)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected extracted file1 to exist")

	extractedFile2 := filepath.Join(tempDir, "file2.txt")
	exists, err = PathExists(extractedFile2)
	assert.NoError(t, err)
	assert.True(t, exists, "Expected extracted file2 to exist")
}

// Helper function to create a temporary zip file for testing.
func createTempZipFile(zipFile string) error {
	// Create a new zip file
	file, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Create a file1 inside the zip file
	file1, err := zipWriter.Create("file1.txt")
	if err != nil {
		return err
	}
	_, err = file1.Write([]byte("Test data for file1"))
	if err != nil {
		return err
	}

	// Create a file2 inside the zip file
	file2, err := zipWriter.Create("file2.txt")
	if err != nil {
		return err
	}
	_, err = file2.Write([]byte("Test data for file2"))
	if err != nil {
		return err
	}

	return nil
}
func TestArrayContains(t *testing.T) {
	var testCases = []struct {
		name        string
		array       []string
		nameToCheck string
		expected    bool
	}{
		{
			name:        "Array contains the name",
			array:       []string{"apple", "banana", "cherry"},
			nameToCheck: "banana",
			expected:    true,
		},
		{
			name:        "Array does not contain the name",
			array:       []string{"apple", "banana", "cherry"},
			nameToCheck: "orange",
			expected:    false,
		},
		{
			name:        "Empty array",
			array:       []string{},
			nameToCheck: "apple",
			expected:    false,
		},
	}

	for _, test := range testCases {
		result := ArrayContains(test.array, test.nameToCheck)
		assert.Equal(t, test.expected, result, test.name)
	}
}
