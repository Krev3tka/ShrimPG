package validator

import "testing"

func TestIsYourPasswordCool(t *testing.T) {
	coolPassword := "repeatingCharacters322!!"
	notCoolPassword := "qwerty123"

	if _, ok := IsYourPasswordCool(coolPassword); !ok {
		t.Error("Test 1 failed")
	}

	if _, ok := IsYourPasswordCool(notCoolPassword); ok {
		t.Error("Test 2 failed")
	}
}
