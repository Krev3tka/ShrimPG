package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/argon2"
)

type Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func DeriveKey(password string, salt []byte, p *Params) ([]byte, error) {

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

func Encrypt(plaintext []byte, password string, p *Params) ([]byte, error) {
	salt, err := GenerateRandomBytes(p.SaltLength)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(password))
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

func Decrypt(ciphertext []byte, password string, p *Params) ([]byte, error) {
	block, err := aes.NewCipher([]byte(password))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceStart := p.SaltLength
	nonceEnd := nonceStart + uint32(gcm.NonceSize())
	nonce := ciphertext[nonceStart:nonceEnd]
	actualCiphertext := ciphertext[nonceEnd:]

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
