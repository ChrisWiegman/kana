package settings

import "fmt"

func (s *Settings) GetURL() string {
	return fmt.Sprintf("%s://%s", s.GetProtocol(), s.GetDomain())
}

func (s *Settings) GetDomain() string {
	return fmt.Sprintf("%s.%s", s.Get("name"), s.Get("domain"))
}

func (s *Settings) GetProtocol() string {
	if s.GetBool("ssl") {
		return "https"
	}

	return "http"
}
