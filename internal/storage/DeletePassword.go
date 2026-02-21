package storage

import (
	"awesomeProject3/internal/model"
	"encoding/json"
	"fmt"
	"os"
)

func Delete(serviceName string) {
	allEntries := make(model.Entry)

	data, err := os.ReadFile(FileName)
	if err != nil {
		fmt.Println("Error: could not read file")
		return
	}
	json.Unmarshal(data, &allEntries)

	if _, ok := allEntries[serviceName]; ok {
		delete(allEntries, serviceName)
		fmt.Println("Password deleted successfully")
	} else {
		fmt.Println("We didn't found the password of the service")
	}

	updatedJson, _ := json.MarshalIndent(allEntries, "", " ")
	os.WriteFile(FileName, updatedJson, 0644)
}
