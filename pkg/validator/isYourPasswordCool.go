package validator

import (
	"fmt"
	"unicode"
)

func IsYourPasswordCool(passwd string) (bool, error) {
	if len(passwd) < 12 {
		return false, fmt.Errorf("password is too short. The minimum required password length is 12 characters, your password has %d symbols", len(passwd))
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

	if !hasDigit {
		return false, fmt.Errorf("password has no digits")
	}
	if !hasUpper {
		return false, fmt.Errorf("password has no uppercase letters")
	}
	if !hasSpecial {
		return false, fmt.Errorf("password has no special symbols")
	}

	return true, nil
}
