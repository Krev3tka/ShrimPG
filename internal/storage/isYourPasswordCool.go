package storage

import "unicode"

func IsYourPasswordCool(passwd string) (string, bool) {
	if len(passwd) < 12 {
		return "You should think harder and come up with a new password", false
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
		return "You need to think about adding uppercase letters, digits, or special characters", false
	}

	return "Your password is cool", true
}
