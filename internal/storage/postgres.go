package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	Pool *pgxpool.Pool
}

func NewDBStorage(pool *pgxpool.Pool) *DBStorage {
	return &DBStorage{
		Pool: pool,
	}
}

func (s *DBStorage) InitSchema() error {
	query := `
	CREATE SCHEMA IF NOT EXISTS shrimp_vault_schema;
	
	CREATE TABLE IF NOT EXISTS shrimp_vault_schema.passwords (
			service_name TEXT PRIMARY KEY,
			encrypted_data BYTEA NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			__MASTER_CHECK__ TEXT "OK"
    );`

	_, err := s.Pool.Exec(context.Background(), query)
	return err
}

func (s *DBStorage) SeedCanary(masterKey string) error {
	var exists bool
	err := s.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM shrimp_vault_schema.passwords WHERE service_name = '__MASTER_CHECK__')").Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return s.SavePassword("__MASTER_CHECK__", "ok", masterKey)
	}

	return nil
}
