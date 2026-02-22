package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Krev3tka/ShrimPG/internal/model"
)

const FileName = "passwords.json"

func Save(Entry model.Entry, masterPassword string) {
	jsonData, err := json.Marshal(Entry)
	if err != nil {
		fmt.Printf("Shrimp is caught and cooked because of: %v", err)
		return
	}

	key := deriveKey(masterPassword)

	encryptedData, err := Encrypt(jsonData, key)
	if err != nil {
		fmt.Printf("Shrimp is caught and cooked because of: %v", err)
		return
	}

	err = os.WriteFile(FileName, encryptedData, 0644)
	if err != nil {
		fmt.Printf("Shrimp is caught and cooked because of: %v", err)
		return
	}
}

func Load(masterPassword string) (model.Entry, error) {
	_, err := os.Stat(FileName)
	if os.IsNotExist(err) {
		fmt.Println("No vault found. Creating a new one for you, shrimp!")
		return model.Entry{}, nil
	}

	encryptedData, err := os.ReadFile(FileName)
	if err != nil {
		fmt.Printf("Shrimp is caught and cooked because of: %v", err)
		return nil, err
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
