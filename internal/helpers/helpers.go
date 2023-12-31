package helpers

import (
	"strings"
)

// IsValidString Checks a given string against an array of valid values and returns true/false as appropriate.
func IsValidString(stringToCheck string, validStrings []string) bool {
	for _, validString := range validStrings {
		if validString == stringToCheck {
			return true
		}
	}

	return false
}

// sanitizeSiteName Returns the site name, properly sanitized for use.
func SanitizeSiteName(rawSiteName string) string {
	siteName := strings.TrimSpace(rawSiteName)
	siteName = strings.ToLower(siteName)
	siteName = strings.ReplaceAll(siteName, " ", "-")
	siteName = strings.ReplaceAll(siteName, "_", "-")
	return strings.ToValidUTF8(siteName, "")
}
