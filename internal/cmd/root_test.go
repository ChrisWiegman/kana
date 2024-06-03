package cmd

import (
	"testing"

	"github.com/ChrisWiegman/kana/tests"
)

func TestRoot(t *testing.T) {
	testCases := []tests.Test{
		{
			Description: "run the kana root command without further input",
			Command:     []string{}},
	}

	tests.RunSnapshotTest(testCases, t)
}
