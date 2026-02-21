package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintList(filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("we couldn't open file because of: %w", err)
		return
	}

	var passwdMap map[string]string
	json.Unmarshal(file, &passwdMap)

	fmt.Println("List of your passwords: ")
	for key, value := range passwdMap {
		fmt.Printf("%s: %s\n", key, value)
	}

}
