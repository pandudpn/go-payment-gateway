package validation

import (
	"regexp"
)

// PhoneNumber validation for indonesian number
// phoneNumber will start with 62 and 11 - 13 digits
func (v *validation) PhoneNumber(phone string) bool {
	if len(phone) < 2 || len(phone) > 13 {
		return false
	}

	// default regex is Indonesia
	var phoneRegexString = "^[62]+[0-9]"

	// when country is from PH
	if v.c == PH {
		phoneRegexString = "^[63]+[0-9]"
	}
	phoneRegex := regexp.MustCompile(phoneRegexString)

	return phoneRegex.MatchString(phone)
}
