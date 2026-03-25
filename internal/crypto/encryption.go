package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

func DeriveKey(password string, salt []byte, p *Argon2Params) ([]byte, error) {

	hash := argon2.Key([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	return hash, nil
}

func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func Encrypt(plaintext []byte, password string, p *Argon2Params) ([]byte, error) {
	block, err := aes.NewCipher([]byte(password))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := GenerateRandomBytes(uint32(gcm.NonceSize()))
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, password string, p *Argon2Params) ([]byte, error) {
	block, err := aes.NewCipher([]byte(password))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce := ciphertext[:nonceSize]
	actualCiphertext := ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
