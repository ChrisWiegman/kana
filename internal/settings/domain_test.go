package settings

import "testing"

func TestGetURL(t *testing.T) {
	options := defaultOptions
	site := Site{
		Name: "example",
	}
	constants := Constants{
		Domain: "com",
	}

	s := &Settings{
		settings:  options,
		constants: constants,
		site:      site,
	}

	expectedURL := "http://example.com"
	actualURL := s.GetURL()

	if actualURL != expectedURL {
		t.Errorf("Expected URL: %s, but got: %s", expectedURL, actualURL)
	}

	s.settings.SSL = true
	expectedURL = "https://example.com"
	actualURL = s.GetURL()

	if actualURL != expectedURL {
		t.Errorf("Expected URL: %s, but got: %s", expectedURL, actualURL)
	}
}

func TestGetDomain(t *testing.T) {
	options := defaultOptions
	site := Site{
		Name: "example",
	}
	constants := Constants{
		Domain: "com",
	}

	s := &Settings{
		settings:  options,
		constants: constants,
		site:      site,
	}

	expectedDomain := "example.com"
	actualDomain := s.GetDomain()
	if actualDomain != expectedDomain {
		t.Errorf("Expected domain: %s, but got: %s", expectedDomain, actualDomain)
	}
}

func TestGetProtocol(t *testing.T) {
	options := defaultOptions
	site := Site{
		Name: "example",
	}
	constants := Constants{
		Domain: "com",
	}

	s := &Settings{
		settings:  options,
		constants: constants,
		site:      site,
	}

	expectedProtocol := "http"
	actualProtocol := s.GetProtocol()
	if actualProtocol != expectedProtocol {
		t.Errorf("Expected protocol: %s, but got: %s", expectedProtocol, actualProtocol)
	}

	s.settings.SSL = true
	expectedProtocol = "https"

	actualProtocol = s.GetProtocol()
	if actualProtocol != expectedProtocol {
		t.Errorf("Expected protocol: %s, but got: %s", expectedProtocol, actualProtocol)
	}
}
