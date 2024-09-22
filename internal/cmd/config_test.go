package cmd

import (
	"testing"

	"github.com/ChrisWiegman/kana-wordpress/tests"
)

func TestConfig(t *testing.T) {
	testCases := []tests.Test{
		{
			Description: "Test the default config command",
			Command:     []string{"config"}},
		{
			Description: "Test the config command with json output",
			Command:     []string{"config", "--output-json"}},
		{
			Description: "Retrieve the PHP value from the config command",
			Command:     []string{"config", "php"}},
		{
			Description: "Retrieve the PHP value from the config command with json output",
			Command:     []string{"config", "php", "--output-json"}},
	}

	tests.RunCommandTest(testCases, t)
}
