package rpc

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/Krev3tka/ShrimPG/internal/crypto"
)

type Service struct {
	client CryptoServiceClient
}

func NewService(client CryptoServiceClient) *Service {
	return &Service{client: client}
}

func (s *Service) GenerateRandomBytes(ctx context.Context, n uint32) ([]byte, error) {
	resp, err := s.client.GenerateRandomBytes(ctx, &GenerateRandomBytesRequest{N: n})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (s *Service) DeriveKey(ctx context.Context, password string, salt []byte) ([]byte, error) {
	resp, err := s.client.DeriveKey(ctx, &DeriveKeyRequest{Password: password, Salt: salt})
	if err != nil {
		return nil, err
	}
	return resp.Key, nil
}

func (s *Service) Encrypt(ctx context.Context, plaintext []byte, key []byte) ([]byte, error) {
	resp, err := s.client.EncryptWithKey(ctx, &EncryptWithKeyRequest{Key: key, Plaintext: plaintext})
	if err != nil {
		return nil, err
	}
	return append(resp.Nonce, resp.Ciphertext...), nil
}

func (s *Service) Decrypt(ctx context.Context, ciphertext []byte, key []byte) ([]byte, error) {
	if len(ciphertext) < 12 {
		return nil, fmt.Errorf("ciphertext too short")
	}
	resp, err := s.client.DecryptWithKey(ctx, &DecryptWithKeyRequest{Key: key, Nonce: ciphertext[:12], Ciphertext: ciphertext[12:]})
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}

func (s *Service) EncryptWithMasterPassword(ctx context.Context, plaintext []byte, masterKey string) ([]byte, error) {
	salt, err := s.GenerateRandomBytes(ctx, crypto.DefaultParams.SaltLength)
	if err != nil {
		return nil, err
	}
	key, err := s.DeriveKey(ctx, masterKey, salt)
	if err != nil {
		return nil, err
	}
	ciphertext, err := s.Encrypt(ctx, plaintext, key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, len(salt)+len(ciphertext))
	out = append(out, salt...)
	out = append(out, ciphertext...)
	return out, nil
}

func (s *Service) DecryptWithMasterPassword(ctx context.Context, encrypted []byte, masterKey string) ([]byte, error) {
	if len(encrypted) < int(crypto.DefaultParams.SaltLength)+12 {
		return nil, fmt.Errorf("encrypted data too short")
	}
	salt := encrypted[:crypto.DefaultParams.SaltLength]
	ciphertext := encrypted[crypto.DefaultParams.SaltLength:]
	key, err := s.DeriveKey(ctx, masterKey, salt)
	if err != nil {
		return nil, err
	}
	return s.Decrypt(ctx, ciphertext, key)
}

func EncodeLenPrefix(data []byte) []byte {
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)
	return buf
}
