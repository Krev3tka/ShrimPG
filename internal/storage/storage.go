package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Krev3tka/ShrimPG/internal/model"
)

const FileName = "passwords.json"

func Save(Entry model.Entry, masterPassword string) error {
	jsonData, err := json.Marshal(Entry)
	if err != nil {
		return fmt.Errorf("Shrimp is caught and cooked because of: %w", err)
	}

	key := deriveKey(masterPassword)

	encryptedData, err := Encrypt(jsonData, key)
	if err != nil {
		return fmt.Errorf("Shrimp is caught and cooked because of: %w", err)
	}

	err = os.WriteFile(FileName, encryptedData, 0644)
	if err != nil {
		return fmt.Errorf("Shrimp is caught and cooked because of: %v", err)
	}

	return nil
}

func Load(masterPassword string) (model.Entry, error) {
	_, err := os.Stat(FileName)
	if os.IsNotExist(err) {
		fmt.Println("No vault found. Creating a new one for you, shrimp!")
		return make(model.Entry), nil
	}

	encryptedData, err := os.ReadFile(FileName)
	if err != nil {
		return nil, fmt.Errorf("Shrimp is caught and cooked because of: %w", err)
	}

	key := deriveKey(masterPassword)

	jsonData, err := Decrypt(encryptedData, key)
	if err != nil {
		return model.Entry{}, fmt.Errorf("wrong master password or corrupted file.")
	}

	var entry model.Entry

	err = json.Unmarshal(jsonData, &entry)
	if err != nil {
		return nil, err
	}

	return entry, nil

}

func Delete(serviceName string, masterPassword string) error {
	vault, err := Load(masterPassword)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	if _, ok := vault[serviceName]; ok {
		delete(vault, serviceName)
		fmt.Println("Password deleted successfully")
	} else {
		fmt.Println("We didn't found the password of the service")
	}

	err = Save(vault, masterPassword)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	return nil
}
