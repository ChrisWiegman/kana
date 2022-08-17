package appConfig

func CheckString(stringToCheck string, validStrings []string) bool {

	for _, validString := range validStrings {
		if validString == stringToCheck {
			return true
		}
	}

	return false

}
