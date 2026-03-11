package crypto

import "testing"

func TestDeriveKey(t *testing.T) {
	p := &Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  12,
		KeyLength:   16,
	}
	password := "pass_word"
	data := []byte("secret_data")
	encrypted, err := Encrypt(data, password, p)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := Decrypt(encrypted, password, p)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != string(decrypted) {
		t.Errorf("Test failed")
	}

}
