package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"golang.org/x/crypto/argon2"
	"io"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func deriveKey(password string, salt []byte, p *params) ([]byte, error) {

	hash := argon2.Key([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return hash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func Encrypt(plaintext []byte, password string, p *params) ([]byte, error) {
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return nil, err
	}

	key, err := deriveKey(password, salt, p)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	result := append(salt, gcm.Seal(nonce, nonce, plaintext, nil)...)
	return result, nil
}

func Decrypt(ciphertext []byte, password string, p *params) ([]byte, error) {
	salt := ciphertext[:p.saltLength]

	key, err := deriveKey(password, salt, p)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceStart := p.saltLength
	nonceEnd := nonceStart + uint32(gcm.NonceSize())
	nonce := ciphertext[nonceStart:nonceEnd]
	actualCiphertext := ciphertext[nonceEnd:]

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
