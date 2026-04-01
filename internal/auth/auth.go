// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package auth

// Package auth provides core utilities for session management and database state verification.
import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// IsDatabaseEmpty reports whether the passwords table contains no records.
// It is used to determine the initial state of the vault during CLI or server startup.
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

// GenerateRandomToken creates a cryptographically secure session token.
//
// The generation process follows these steps:
//  1. A 16-byte buffer is allocated.
//  2. The buffer is filled with random data using [crypto/rand].
//  3. The data is encoded into a hexadecimal string.
//
// It returns a 32-character string or an error if the entropy source fails.
func GenerateRandomToken() (string, error) {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
