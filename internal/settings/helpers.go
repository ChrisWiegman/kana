package settings

import (
	"strings"
)

// isValidString Checks a given string against an array of valid values and returns true/false as appropriate
func isValidString(stringToCheck string, validStrings []string) bool {
	for _, validString := range validStrings {
		if validString == stringToCheck {
			return true
		}
	}

	return false
}

// sanitizeSiteName Returns the site name, properly sanitized for use.
func sanitizeSiteName(rawSiteName string) string {
	siteName := strings.TrimSpace(rawSiteName)
	siteName = strings.ToLower(siteName)
	siteName = strings.ReplaceAll(siteName, " ", "-")
	return strings.ToValidUTF8(siteName, "")
}
