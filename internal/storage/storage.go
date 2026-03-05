package storage

import (
	"context"
	"fmt"

	"github.com/Krev3tka/ShrimPG/internal/model"
)

const FileName = "passwords.json"

func (s *DBStorage) SavePassword(service string, password string, masterKey string) error {
	key := deriveKey(masterKey)
	query := `
		INSERT INTO shrimp_vault_schema.passwords (service_name, encrypted_data)
		VALUES ($1, $2)
	`

	fmt.Printf("DEBUG: Executing query for service: %s\n", service)

	encrypted, err := Encrypt([]byte(password), key)

	_, err = s.Pool.Exec(context.Background(), query, service, encrypted)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) GetPassword(serviceName string, masterKey string) ([]byte, error) {
	key := deriveKey(masterKey)
	query := `SELECT encrypted_data FROM shrimp_vault_schema.passwords WHERE service_name = $1`

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
	query := `DELETE FROM shrimp_vault_schema.passwords WHERE service_name = $1`

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
	query := `SELECT service_name, encrypted_data FROM shrimp_vault_schema.passwords`

	rows, err := s.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

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
