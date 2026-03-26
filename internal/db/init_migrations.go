package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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

			_, err = dbPool.Exec(context.Background(), string(content))
			if err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", file.Name(), err)
			}
			slog.Info("applied migration", "file", file.Name())
		}
	}

	return nil
}
