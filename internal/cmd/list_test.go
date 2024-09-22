package cmd

import (
	"testing"

	"github.com/ChrisWiegman/kana-dev/tests"
)

func TestList(t *testing.T) {
	testCases := []tests.Test{
		{
			Description: "Test the default list command",
			Command:     []string{"list"}},
		{
			Description: "Test the list command with json output",
			Command:     []string{"list", "--output-json"}},
	}

	tests.RunCommandTest(testCases, t)
}
