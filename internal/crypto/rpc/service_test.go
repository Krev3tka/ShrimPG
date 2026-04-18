package rpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

type fakeClient struct{}

func (f fakeClient) GenerateRandomBytes(ctx context.Context, in *GenerateRandomBytesRequest, opts ...grpc.CallOption) (*GenerateRandomBytesResponse, error) {
	return &GenerateRandomBytesResponse{Data: []byte{1, 2, 3}}, nil
}

func (f fakeClient) DeriveKey(ctx context.Context, in *DeriveKeyRequest, opts ...grpc.CallOption) (*DeriveKeyResponse, error) {
	return &DeriveKeyResponse{Key: []byte("derived-key-derived-key-derived-key!!")[:32]}, nil
}

func (f fakeClient) EncryptWithKey(ctx context.Context, in *EncryptWithKeyRequest, opts ...grpc.CallOption) (*EncryptWithKeyResponse, error) {
	return &EncryptWithKeyResponse{Nonce: []byte("123456789012"), Ciphertext: []byte{9, 9, 9}}, nil
}

func (f fakeClient) DecryptWithKey(ctx context.Context, in *DecryptWithKeyRequest, opts ...grpc.CallOption) (*DecryptWithKeyResponse, error) {
	return &DecryptWithKeyResponse{Plaintext: []byte("secret")}, nil
}

func TestServiceRoundTripHelpers(t *testing.T) {
	service := NewService(fakeClient{})

	gotBytes, err := service.GenerateRandomBytes(context.Background(), 3)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotBytes) != string([]byte{1, 2, 3}) {
		t.Fatalf("unexpected bytes: %v", gotBytes)
	}

	key, err := service.DeriveKey(context.Background(), "password", []byte("salt-salt-salt-"))
	if err != nil {
		t.Fatal(err)
	}
	if len(key) != 32 {
		t.Fatalf("unexpected key length: %d", len(key))
	}

	enc, err := service.Encrypt(context.Background(), []byte("hello"), key)
	if err != nil {
		t.Fatal(err)
	}
	if len(enc) == 0 {
		t.Fatal("empty encrypted output")
	}

	plain, err := service.Decrypt(context.Background(), enc, key)
	if err != nil {
		t.Fatal(err)
	}
	if string(plain) != "secret" {
		t.Fatalf("unexpected plaintext: %s", plain)
	}
}
