package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Connect() (*pgxpool.Pool, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, fmt.Errorf("error while loading .env file: %w", err)
	}

	connStr := os.Getenv("CONN_STR")

	dbPoll, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return dbPoll, nil
}
