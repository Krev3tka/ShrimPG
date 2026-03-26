package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func IsDatabaseEmpty(dbPool *pgxpool.Pool) bool {
	query := "SELECT count(*) FROM passwords"
	ctx := context.Background()

	var count int
	err := dbPool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		slog.Error("Failed to get count of passwords", "details", err)
	}

	return count == 0
}

func GenerateRandomToken() (string, error) {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
