package db

import (
	"context"
	"fmt"

	"github.com/Krev3tka/ShrimPG/internal/crypto"
	"github.com/Krev3tka/ShrimPG/internal/model"
)

func (s *DBStorage) SavePassword(userID int, service string, passwd string, masterKey string, p *crypto.Params) error {
	salt, err := crypto.GenerateRandomBytes(p.SaltLength)
	if err != nil {
		return err
	}

	key, err := crypto.DeriveKey(masterKey, salt, p)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	encrypted, err := crypto.Encrypt([]byte(passwd), string(key), p)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	query := "INSERT INTO passwords (user_id, service, encrypted_data) VALUES ($1, $2, $3)"
	_, err = s.Pool.Exec(context.Background(), query, userID, service, encrypted)
	return err
}

func (s *DBStorage) GetPassword(serviceName string, masterKey string, p *crypto.Params) ([]byte, error) {
	salt, err := crypto.GenerateRandomBytes(p.SaltLength)
	if err != nil {
		return nil, err
	}

	key, err := crypto.DeriveKey(masterKey, salt, p)
	if err != nil {
		return nil, err
	}

	query := "SELECT encrypted_data FROM passwords WHERE service = $1"

	var encryptedData []byte
	err = s.Pool.QueryRow(context.Background(), query, serviceName).Scan(&encryptedData)
	if err != nil {
		return nil, err
	}

	decryptedData, err := crypto.Decrypt(encryptedData, string(key), p)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func (s *DBStorage) DeletePassword(service string) error {
	query := "DELETE FROM passwords WHERE service = $1"

	result, err := s.Pool.Exec(context.Background(), query, service)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("service not found")
	}
	return nil
}

func (s *DBStorage) GetAllPasswords(masterKey string, p *crypto.Params) (model.Entry, error) {
	salt, err := crypto.GenerateRandomBytes(p.SaltLength)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}
	key, err := crypto.DeriveKey(masterKey, salt, p)
	if err != nil {
		return nil, err
	}
	query := "SELECT service, encrypted_data FROM passwords"

	rows, err := s.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entry := make(model.Entry)

	for rows.Next() {
		var serviceName string
		var passwd []byte

		err := rows.Scan(&serviceName, &passwd)
		if err != nil {
			return nil, err
		}

		decryptedData, err := crypto.Decrypt(passwd, string(key), p)
		if err != nil {
			return nil, err
		}

		entry[serviceName] = string(decryptedData)
	}

	return entry, nil
}
