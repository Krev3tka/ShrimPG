// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Krev3tka/ShrimPG/internal/crypto"
	"github.com/Krev3tka/ShrimPG/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *DBStorage) VerifyAuthHash(ctx context.Context, username string, authHash string) (int, error) {
	var dbHash []byte
	var userID int

	query := "SELECT id, master_hash FROM users WHERE name = $1"
	err := s.Pool.QueryRow(ctx, query, username).Scan(&userID, &dbHash)
	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(dbHash, []byte(authHash))
	if err != nil {
		return 0, fmt.Errorf("invalid credentials")
	}

	return userID, nil
}

func (s *DBStorage) CreateUser(ctx context.Context, username string, authHash string) (int, error) {
	clientSalt, _ := crypto.GenerateRandomBytes(16)

	serverHash, err := bcrypt.GenerateFromPassword([]byte(authHash), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var id int
	query := "INSERT INTO users (name, master_hash, master_salt) VALUES ($1, $2, $3) RETURNING id"
	err = s.Pool.QueryRow(ctx, query, username, serverHash, clientSalt).Scan(&id)

	return id, err

}

func (s *DBStorage) SavePassword(userID int, service []byte, encryptedData []byte) error {
	query := "INSERT INTO passwords (user_id, service, encrypted_data) VALUES ($1, $2, $3)"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.Pool.Exec(ctx, query, userID, service, encryptedData)
	if err != nil {
		return fmt.Errorf("db exec error: %w", err)
	}
	return nil
}

func (s *DBStorage) GetPassword(userID int, serviceName string) ([]byte, error) {
	query := "SELECT encrypted_data FROM passwords WHERE service = $1 AND user_id = $2"
	var encryptedData []byte

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Pool.QueryRow(ctx, query, serviceName, userID).Scan(&encryptedData)
	if err != nil {
		return nil, fmt.Errorf("db: get password: %w", err)
	}

	return encryptedData, nil
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

func (s *DBStorage) GetAllPasswords(userID int) (model.Entry, error) {
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

		entry[serviceName] = encryptedData
	}

	return entry, nil
}
