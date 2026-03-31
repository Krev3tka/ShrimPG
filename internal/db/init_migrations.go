// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitMigrations(dbPool *pgxpool.Pool, migrationsDir string) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			slog.Info("applying migration", "file", file.Name())

			content, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
			}

			commands := strings.Split(string(content), ";")
			for _, cmd := range commands {
				trimmedCmd := strings.TrimSpace(cmd)
				if trimmedCmd == "" {
					continue
				}

				execCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_, err = dbPool.Exec(execCtx, trimmedCmd)
				cancel()

				if err != nil {
					return fmt.Errorf("failed to apply command in %s: %w", file.Name(), err)
				}
			}

			slog.Info("applied migration", "file", file.Name())
		}
	}

	return nil
}
