package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Krev3tka/ShrimPG/internal/api"
	"github.com/Krev3tka/ShrimPG/internal/auth"
	"github.com/Krev3tka/ShrimPG/internal/db"
	"github.com/Krev3tka/ShrimPG/internal/storage"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	masterKey := auth.GetMasterPassword()

	dbPool, err := db.Connect()
	if err != nil {
		slog.Error("Database connection failed", "details", err)
		return
	}
	defer dbPool.Close()

	slog.Info("database connection established", "address", "localhost:5432")

	vault := storage.NewDBStorage(dbPool)

	handler := api.NewHandler(vault, masterKey)

	http.HandleFunc("/passwords/create", handler.AuthMiddleware(handler.CreatePasswordRequest))
	http.HandleFunc("/passwords/get", handler.AuthMiddleware(handler.GetPasswordRequest))
	http.HandleFunc("/passwords/delete", handler.AuthMiddleware(handler.DeletePasswordRequest))

	slog.Info("Server is starting on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Server crashed", "error", err)
		os.Exit(1)
	}
}
