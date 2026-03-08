package storage

import (
	"context"
	"fmt"

	"github.com/Krev3tka/ShrimPG/internal/model"
)

func (s *DBStorage) SavePassword(userID int, service string, passwd string, masterKey string) error {
	key := deriveKey(masterKey)

	encrypted, err := Encrypt([]byte(passwd), key)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	query := "INSERT INTO passwords (user_id, service, encrypted_data) VALUES ($1, $2, $3)"
	_, err = s.Pool.Exec(context.Background(), query, userID, service, encrypted)
	return err
}

func (s *DBStorage) GetPassword(serviceName string, masterKey string) ([]byte, error) {
	key := deriveKey(masterKey)
	query := "SELECT encrypted_data FROM passwords WHERE service = $1"

	var encryptedData []byte
	err := s.Pool.QueryRow(context.Background(), query, serviceName).Scan(&encryptedData)
	if err != nil {
		return nil, err
	}

	decryptedData, err := Decrypt(encryptedData, key)
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

func (s *DBStorage) GetAllPasswords(masterKey string) (model.Entry, error) {
	key := deriveKey(masterKey)
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

		decryptedData, err := Decrypt(passwd, key)
		if err != nil {
			return nil, err
		}

		entry[serviceName] = string(decryptedData)
	}

	return entry, nil
}
