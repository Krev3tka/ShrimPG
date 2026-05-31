package validator

import "testing"

func TestValidatePassword(t *testing.T) {
	coolPassword := "repeatingCharacters322!!"
	notCoolPassword := "qwerty123"

	if ok, _ := ValidatePassword(coolPassword); !ok {
		t.Error("Test 1 failed")
	}

	if ok, _ := ValidatePassword(notCoolPassword); ok {
		t.Error("Test 2 failed")
	}
}
