package helpers

import (
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
