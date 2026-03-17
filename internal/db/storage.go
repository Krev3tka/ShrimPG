package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Krev3tka/ShrimPG/internal/crypto"
	"github.com/Krev3tka/ShrimPG/internal/model"
	"github.com/Krev3tka/ShrimPG/pkg/validator"
)

func (s *DBStorage) SavePassword(userID int, service string, passwd string, masterKey string) error {
	if ok, err := validator.IsYourPasswordCool(passwd); !ok {
		return fmt.Errorf("your password isn't safe yet: %w", err)
	}
	salt, err := crypto.GenerateRandomBytes(s.Config.params.SaltLength)
	if err != nil {
		return fmt.Errorf("db: save password: %w", err)
	}

	key, err := crypto.DeriveKey(masterKey, salt, s.Config.params)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	encrypted, err := crypto.Encrypt([]byte(passwd), string(key), s.Config.params)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	query := "INSERT INTO passwords (user_id, salt, service, encrypted_data) VALUES ($1, $2, $3, $4)"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = s.Pool.Exec(ctx, query, userID, salt, service, encrypted)
	if err != nil {
		return fmt.Errorf("db exec error: %w", err)
	}
	return nil
}

func (s *DBStorage) GetPassword(serviceName string, masterKey string) ([]byte, error) {
	query := "SELECT encrypted_data, salt FROM passwords WHERE service = $1"

	var encryptedData []byte
	var salt []byte
	err := s.Pool.QueryRow(context.Background(), query, serviceName).Scan(&encryptedData, &salt)
	if err != nil {
		return nil, fmt.Errorf("db: get password: %w", err)
	}

	key, err := crypto.DeriveKey(masterKey, salt, s.Config.params)
	if err != nil {
		return nil, fmt.Errorf("db: get password: %w", err)
	}

	decryptedData, err := crypto.Decrypt(encryptedData, string(key), s.Config.params)
	if err != nil {
		return nil, fmt.Errorf("db: get password: %w", err)
	}

	return decryptedData, nil
}

func (s *DBStorage) DeletePassword(service string) error {
	query := "DELETE FROM passwords WHERE service = $1"

	result, err := s.Pool.Exec(context.Background(), query, service)
	if err != nil {
		return fmt.Errorf("db: delete password: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("service not found")
	}
	return nil
}

func (s *DBStorage) GetAllPasswords(masterKey string) (model.Entry, error) {
	query := "SELECT salt, service, encrypted_data FROM passwords"

	rows, err := s.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("db: get all passwords: %w", err)
	}
	defer rows.Close()

	entry := make(model.Entry)

	for rows.Next() {
		var serviceName string
		var passwd []byte
		var rowSalt []byte

		err := rows.Scan(&rowSalt, &serviceName, &passwd)
		if err != nil {
			return nil, fmt.Errorf("db: get all passwords: %w", err)
		}

		key, err := crypto.DeriveKey(masterKey, rowSalt, s.Config.params)
		if err != nil {
			return nil, fmt.Errorf("db: get all passwords: %w", err)
		}

		decryptedData, err := crypto.Decrypt(passwd, string(key), s.Config.params)
		if err != nil {
			return nil, fmt.Errorf("db: get all passwords: %w", err)
		}

		entry[serviceName] = string(decryptedData)
	}

	return entry, nil
}
