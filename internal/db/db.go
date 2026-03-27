// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Connect() (*pgxpool.Pool, error) {
	paths := []string{".env", "../../.env"}
	for _, path := range paths {
		if err := godotenv.Load(path); err != nil {
			break
		}
	}

	connStr := os.Getenv("CONN_STR")
	if connStr == "" {
		return nil, fmt.Errorf("CONN_STR environment variable is not set")
	}

	dbPoll, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return dbPoll, nil
}
