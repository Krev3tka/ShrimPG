package storage

import (
	"awesomeProject3/internal/model"
	"encoding/json"
	"fmt"
	"os"
)

const FileName = "passwords.json"

func Save(Entry model.Entry) {
	allEntries := make(model.Entry)

	content, err := os.ReadFile(FileName)
	if err == nil && len(content) > 0 {
		json.Unmarshal(content, &allEntries)
	}

	for key, value := range Entry {
		allEntries[key] = value
	}

	jsonData, err := json.MarshalIndent(allEntries, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}

	err = os.WriteFile(FileName, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}
