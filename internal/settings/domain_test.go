package settings

import "testing"

func TestGetURL(t *testing.T) {
	var tests = []struct {
		name          string
		expectedURL   string
		settingsArray []Setting
	}{
		{
			name:        "No settings",
			expectedURL: "http://.sites.kana.sh",
		},
		{
			name:        "Name is set, but nothing else",
			expectedURL: "http://test.sites.kana.sh",
			settingsArray: []Setting{
				{
					name:         "name",
					currentValue: "test",
				},
			},
		},
		{
			name:        "Name is set, and SSL is true",
			expectedURL: "https://test.sites.kana.sh",
			settingsArray: []Setting{
				{
					name:         "name",
					currentValue: "test",
				},
				{
					name:         "ssl",
					currentValue: "true",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := new(Settings)
			s.settings = test.settingsArray

			actualURL := s.GetURL()

			if actualURL != test.expectedURL {
				t.Errorf("Unexpected URL. Expected: %s, Got: %s", test.expectedURL, actualURL)
			}
		})
	}
}

func TestGetDomain(t *testing.T) {
	settingsArray := []Setting{
		{
			name:         "name",
			currentValue: "test",
		},
	}

	s := new(Settings)
	s.settings = settingsArray

	expectedDomain := "test.sites.kana.sh"

	actualDomain := s.GetDomain()
	if actualDomain != expectedDomain {
		t.Errorf("Unexpected domain. Expected: %s, Got: %s", expectedDomain, actualDomain)
	}
}

func TestGetProtocol(t *testing.T) {
	var tests = []struct {
		name             string
		expectedProtocol string
		settingsArray    []Setting
	}{
		{
			name:             "No settings",
			expectedProtocol: "http",
		},
		{
			name:             "SSL is set to true",
			expectedProtocol: "https",
			settingsArray: []Setting{
				{
					name:         "ssl",
					currentValue: "true",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := new(Settings)
			s.settings = test.settingsArray

			actualProtocol := s.GetProtocol()
			if actualProtocol != test.expectedProtocol {
				t.Errorf("Unexpected protocol. Expected: %s, Got: %s", test.expectedProtocol, actualProtocol)
			}
		})
	}
}
