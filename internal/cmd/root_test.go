package cmd

import (
	"testing"

	"github.com/ChrisWiegman/kana-wp/tests"
)

func TestRoot(t *testing.T) {
	testCases := []tests.Test{
		{
			Description: "run the kana root command without further input",
			Command:     []string{}},
	}

	tests.RunCommandTest(testCases, t)
}
