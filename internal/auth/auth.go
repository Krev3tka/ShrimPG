package auth

import (
	"context"
	"log/slog"
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

	var password string

	if strings.TrimSpace(password) == "" {
	}

	return password
}
