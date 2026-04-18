package db

import (
	"context"
	"errors"
	"testing"
)

type fakeEngine struct{}

func (fakeEngine) GenerateRandomBytes(ctx context.Context, n uint32) ([]byte, error) {
	return make([]byte, n), nil
}
func (fakeEngine) DeriveKey(ctx context.Context, password string, salt []byte) ([]byte, error) {
	return []byte("01234567890123456789012345678901"), nil
}
func (fakeEngine) Encrypt(ctx context.Context, plaintext []byte, key []byte) ([]byte, error) {
	return append([]byte("nonce-nonce12"), plaintext...), nil
}
func (fakeEngine) Decrypt(ctx context.Context, ciphertext []byte, key []byte) ([]byte, error) {
	if len(ciphertext) < len("nonce-nonce12") {
		return nil, errors.New("short")
	}
	return ciphertext[len("nonce-nonce12"):], nil
}

func TestDBStorageUsesInjectedCryptoEngine(t *testing.T) {
	storage := NewDBStorageWithCrypto(nil, fakeEngine{})
	if _, err := storage.generateRandomBytes(context.Background(), 4); err != nil {
		t.Fatal(err)
	}
	if _, err := storage.deriveKey(context.Background(), "pw", []byte("salt-salt-salt-1")); err != nil {
		t.Fatal(err)
	}
	enc, err := storage.encrypt(context.Background(), []byte("hello"), []byte("01234567890123456789012345678901"))
	if err != nil {
		t.Fatal(err)
	}
	plain, err := storage.decrypt(context.Background(), enc, []byte("01234567890123456789012345678901"))
	if err != nil {
		t.Fatal(err)
	}
	if string(plain) != "hello" {
		t.Fatalf("unexpected plaintext: %s", plain)
	}
}
