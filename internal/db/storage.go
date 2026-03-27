// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package db

import (
	"context"
	"crypto/subtle"
	"fmt"
	"time"

	"github.com/Krev3tka/ShrimPG/internal/crypto"
	"github.com/Krev3tka/ShrimPG/internal/model"
	"github.com/Krev3tka/ShrimPG/pkg/validator"
)

func (s *DBStorage) VerifyMasterKey(ctx context.Context, username string, masterKey string) (int, []byte, error) {
	var dbHash, salt []byte
	var userID int

	query := "SELECT id, master_hash, master_salt FROM users WHERE name = $1"
	err := s.Pool.QueryRow(ctx, query, username).Scan(&userID, &dbHash, &salt)
	if err != nil {
		return 0, nil, fmt.Errorf("user not found: %w", err)
	}

	key, err := crypto.DeriveKey(masterKey, salt, s.Config.params)
	if err != nil {
		return 0, nil, err
	}

	if ok := subtle.ConstantTimeCompare(key, dbHash); ok != 1 {
		return 0, nil, fmt.Errorf("invalid master password")
	}

	return userID, key, nil
}

func (s *DBStorage) CreateUser(ctx context.Context, username string, masterKey string) (int, error) {
	salt, err := crypto.GenerateRandomBytes(s.Config.params.SaltLength)

	if err != nil {
		return 0, err
	}

	hash, err := crypto.DeriveKey(masterKey, salt, s.Config.params)

	if err != nil {
		return 0, err
	}

	var id int
	query := "INSERT INTO users (name, master_hash, master_salt) VALUES ($1, $2, $3) RETURNING id"

	err = s.Pool.QueryRow(ctx, query, username, hash, salt).Scan(&id)

	return id, err

}

func (s *DBStorage) SavePassword(userID int, service string, passwd string, encryptionKey []byte) error {
	if ok, err := validator.ValidatePassword(passwd); !ok {
		return fmt.Errorf("your password isn't safe yet: %w", err)
	}

	encrypted, err := crypto.Encrypt([]byte(passwd), string(encryptionKey), s.Config.params)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	query := "INSERT INTO passwords (user_id, service, encrypted_data) VALUES ($1, $2, $3)"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = s.Pool.Exec(ctx, query, userID, service, encrypted)
	if err != nil {
		return fmt.Errorf("db exec error: %w", err)
	}
	return nil
}

func (s *DBStorage) GetPassword(userID int, serviceName string, encryptionKey []byte) ([]byte, error) {
	query := "SELECT encrypted_data FROM passwords WHERE service = $1 AND user_id = $2"

	var encryptedData []byte

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Pool.QueryRow(ctx, query, serviceName, userID).Scan(&encryptedData)
	if err != nil {
		return nil, fmt.Errorf("db: get password: %w", err)
	}

	decryptedData, err := crypto.Decrypt(encryptedData, string(encryptionKey), s.Config.params)
	if err != nil {
		return nil, fmt.Errorf("db: get password: %w", err)
	}

	return decryptedData, nil
}

func (s *DBStorage) DeletePassword(userID int, service string) error {
	query := "DELETE FROM passwords WHERE service = $1 AND user_id = $2"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := s.Pool.Exec(ctx, query, service, userID)
	if err != nil {
		return fmt.Errorf("db: delete password: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("service not found")
	}
	return nil
}

func (s *DBStorage) GetAllPasswords(userID int, encryptionKey []byte) (model.Entry, error) {
	query := "SELECT service, encrypted_data FROM passwords WHERE user_id = $1"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("db: get all passwords: %w", err)
	}

	defer rows.Close()

	entry := make(model.Entry)

	for rows.Next() {
		var serviceName string
		var encryptedData []byte

		err := rows.Scan(&serviceName, &encryptedData)
		if err != nil {
			return nil, fmt.Errorf("db: get all passwords: %w", err)
		}

		decryptedData, err := crypto.Decrypt(encryptedData, string(encryptionKey), s.Config.params)
		if err != nil {
			continue
		}

		entry[serviceName] = string(decryptedData)
	}

	return entry, nil
}
