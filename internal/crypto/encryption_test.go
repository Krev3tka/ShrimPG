// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package crypto

import "testing"

func TestDeriveKey(t *testing.T) {
	p := &DefaultParams
	password := "Coolpass_word123!"
	data := []byte("secret_data")

	salt, err := GenerateRandomBytes(p.SaltLength)
	if err != nil {
		t.Fatal(err)
	}

	key, err := DeriveKey(password, salt, p)

	encrypted, err := Encrypt(data, string(key), p)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := Decrypt(encrypted, string(key), p)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != string(decrypted) {
		t.Errorf("Test failed")
	}

}
