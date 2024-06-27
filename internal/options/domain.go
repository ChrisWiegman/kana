package options

import "fmt"

func (s *Settings) GetURL() string {
	return fmt.Sprintf("%s://%s", s.GetProtocol(), s.GetDomain())
}

func (s *Settings) GetDomain() string {
	return fmt.Sprintf("%s.%s", s.Get("name"), domain)
}

func (s *Settings) GetProtocol() string {
	if s.GetBool("ssl") {
		return "https"
	}

	return "http"
}
