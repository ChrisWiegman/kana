package settings

import (
	"testing"
)

func TestSettings_Get(t *testing.T) {
	s := new(Settings)

	for i := range defaults {
		defaults[i].currentValue = defaults[i].defaultValue
		s.settings = append(s.settings, defaults[i])
	}

	for _, tt := range s.settings {
		t.Run(tt.name, func(t *testing.T) {
			got := s.Get(tt.name)
			if got != tt.defaultValue {
				t.Errorf("Got %q, expected %q", got, tt.defaultValue)
			}
		})
	}
}
func TestSettings_GetBool(t *testing.T) {
	s := &Settings{
		settings: []Setting{
			{name: "debug", currentValue: "true"},
			{name: "verbose", currentValue: "false"},
			{name: "type", currentValue: "site"},
		},
	}

	tests := []struct {
		name     string
		expected bool
	}{
		{"debug", true},
		{"verbose", false},
		{"nonexistent", false},
		{"type", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.GetBool(tt.name)
			if got != tt.expected {
				t.Errorf("Got %v, expected %v", got, tt.expected)
			}
		})
	}
}
