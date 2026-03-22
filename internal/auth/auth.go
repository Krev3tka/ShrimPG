package auth

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetMasterPassword(dbPool *pgxpool.Pool) string {
	query := "SELECT count(*) FROM passwords"
	ctx := context.Background()

	var count int
	err := dbPool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		slog.Error("Failed to get count of passwords", "details", err)
	}

	if count == 0 {
		fmt.Println("=== FIRST RUN: Create your Master Password ===")
	} else {
		fmt.Println("=== AUTH: Enter Master Password ===")
	}
	fmt.Print("Enter Master Password: ")

	//bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	//fmt.Println()

	var password string
	_, err = fmt.Scanln(&password)

	if err != nil {
		fmt.Println("Error reading password.")
		os.Exit(1)
	}

	//password := string(bytePassword)
	if strings.TrimSpace(password) == "" {
		fmt.Println("Password is empty")
		os.Exit(1)
	}

	return password
}
