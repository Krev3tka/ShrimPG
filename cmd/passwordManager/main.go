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
	"github.com/Krev3tka/ShrimPG/internal/db"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

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

	handler := api.NewHandler(vault, 0)

	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/logout", handler.AuthMiddleware(handler.Logout))
	http.HandleFunc("/passwords/create", handler.AuthMiddleware(handler.CreatePasswordRequest))
	http.HandleFunc("/passwords/get", handler.AuthMiddleware(handler.GetPasswordRequest))
	http.HandleFunc("/passwords/delete", handler.AuthMiddleware(handler.DeletePasswordRequest))
	http.HandleFunc("/passwords/list", handler.AuthMiddleware(handler.GetAllPasswordsRequest))

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	slog.Info("Server is starting.", "port", port)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server crashed unexpectedly", "error", err)
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
