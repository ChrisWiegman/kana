package cmd

import (
	"testing"

	"github.com/ChrisWiegman/kana-wordpress/tests"
)

func TestVersion(t *testing.T) {
	testCases := []tests.Test{
		{
			Description: "Test the version command for appropriate output",
			Command:     []string{"version"}},
		{
			Description: "Test the config command with json output",
			Command:     []string{"version", "--output-json"}},
	}

	tests.RunCommandTest(testCases, t)
}
