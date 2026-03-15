package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Krev3tka/ShrimPG/internal/api"
	"github.com/Krev3tka/ShrimPG/internal/auth"
	"github.com/Krev3tka/ShrimPG/internal/db"
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
	defer func() {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		dbPool.Close()
	}()

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("database connection established", "address", "localhost:5432")

	vault := db.NewDBStorage(dbPool)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := vault.Ping(ctx); err != nil {
		slog.Error("Database is unreachable", "error", err)
		return
	}

	handler := api.NewHandler(vault, masterKey)

	http.HandleFunc("/passwords/create", handler.AuthMiddleware(handler.CreatePasswordRequest))
	http.HandleFunc("/passwords/get", handler.AuthMiddleware(handler.GetPasswordRequest))
	http.HandleFunc("/passwords/delete", handler.AuthMiddleware(handler.DeletePasswordRequest))

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	slog.Info("Server is starting.", "port", port)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("Server crashed", "error", err)
			os.Exit(1)
		}
	}()

	<-exit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		slog.Error("Error: failed to shutdown server correctly", "error", err)
		return
	}
	slog.Info("Server stopped gracefully")
}
