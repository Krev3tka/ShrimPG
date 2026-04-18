// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package db

import (
	"context"
	"fmt"

	"github.com/Krev3tka/ShrimPG/internal/crypto"
	cryptorpc "github.com/Krev3tka/ShrimPG/internal/crypto/rpc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CryptoEngine interface {
	GenerateRandomBytes(ctx context.Context, n uint32) ([]byte, error)
	DeriveKey(ctx context.Context, password string, salt []byte) ([]byte, error)
	Encrypt(ctx context.Context, plaintext []byte, key []byte) ([]byte, error)
	Decrypt(ctx context.Context, ciphertext []byte, key []byte) ([]byte, error)
}

type Config struct {
	params *crypto.Argon2Params
}

type DBStorage struct {
	Pool   *pgxpool.Pool
	Config Config
	Crypto CryptoEngine
}

func NewDBStorage(pool *pgxpool.Pool) *DBStorage {
	return &DBStorage{
		Pool: pool,
		Config: Config{
			params: &crypto.DefaultParams,
		},
	}
}

func NewDBStorageWithCrypto(pool *pgxpool.Pool, engine CryptoEngine) *DBStorage {
	storage := NewDBStorage(pool)
	storage.Crypto = engine
	return storage
}

func NewCryptoEngineFromConn(conn cryptorpc.CryptoServiceClient) CryptoEngine {
	return cryptorpc.NewService(conn)
}

func (s *DBStorage) Ping(ctx context.Context) error {
	query := "SELECT 1"
	_, err := s.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("db: Ping: failed to ping database: %w", err)
	}
	return nil
}
