// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package db

import (
	"context"
	"fmt"

	"github.com/Krev3tka/ShrimPG/internal/crypto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	params *crypto.Argon2Params
}

type DBStorage struct {
	Pool   *pgxpool.Pool
	Config Config
}

func NewDBStorage(pool *pgxpool.Pool) *DBStorage {
	return &DBStorage{
		Pool: pool,
		Config: Config{
			params: &crypto.DefaultParams,
		},
	}
}

func (s *DBStorage) Ping(ctx context.Context) error {
	query := "SELECT 1"
	_, err := s.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("db: Ping: failed to ping database: %w", err)
	}
	return nil
}
