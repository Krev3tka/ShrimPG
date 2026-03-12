package validator

import "unicode"

func IsYourPasswordCool(passwd string) bool {
	if len(passwd) < 12 {
		return false
	}

	var hasDigit, hasUpper, hasSpecial bool

	for _, char := range passwd {
		switch {
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasDigit || !hasSpecial {
		return false
	}

	return true
}
