package cmd

import (
	"testing"

	"github.com/ChrisWiegman/kana-wp/tests"
)

func TestChangelog(t *testing.T) {
	testCases := []tests.Test{
		{
			Description: "Test the default changelog command",
			Command:     []string{"changelog"},
			Output:      "The Kana changelog has been opened in your default browser."},
		{
			Description: "Test the changelog command with json output",
			Command:     []string{"changelog", "--output-json"},
			Output:      `{"Status":"Success","Message":"The Kana changelog has been opened in your default browser."}`},
	}

	tests.RunCommandTest(testCases, t)
}
