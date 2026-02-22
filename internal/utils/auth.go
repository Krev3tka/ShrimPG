package utils

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func GetMasterPassword() string {
	fmt.Print("Enter Master Password: ")

	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	if err != nil {
		fmt.Println("Error reading password.")
		os.Exit(1)
	}

	password := string(bytePassword)
	if strings.TrimSpace(password) == "" {
		fmt.Println("Empty password? Are you serious?\n(It's time to introduce a sobriety test)")
		os.Exit(1)
	}

	return password
}
