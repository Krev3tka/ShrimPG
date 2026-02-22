package storage

import (
	"fmt"

	"github.com/Krev3tka/ShrimPG/internal/model"
)

func PrintVault(Entry model.Entry) {
	fmt.Println("List of your passwords: ")
	for key, value := range Entry {
		fmt.Printf("%s: %s\n", key, value)
	}

}
