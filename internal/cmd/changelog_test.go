package cmd

import (
	"testing"

	"github.com/ChrisWiegman/kana/tests"
)

func TestChangelog(t *testing.T) {
	testCases := []tests.Test{
		{
			Description: "Test the default changelog command",
			Command:     []string{"changelog"}},
		{
			Description: "Test the changelog command with json output",
			Command:     []string{"changelog", "--output-json"}},
	}

	tests.RunSnapshotTest(testCases, t)
}
