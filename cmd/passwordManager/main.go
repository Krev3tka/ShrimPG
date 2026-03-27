// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
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

	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "0.0.0.0:8080"
	}

	u, err := url.Parse(serverAddr)
	if err == nil {
		slog.Info("database connection established", "host", u.Host)
	}
	vault := db.NewDBStorage(dbPool)
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		err = dbPool.Ping(context.Background())
		if err == nil {
			break
		}
		slog.Info("Waiting for database...", "attempt", i+1)
		time.Sleep(2 * time.Second)
	}

	slog.Info("running migrations...")
	if err := db.InitMigrations(dbPool, "./migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		return
	}

	handler := api.NewHandler(vault)

	mux := http.NewServeMux()
	mux.HandleFunc("/register", api.CORSmiddleware(handler.Register))
	mux.HandleFunc("/login", api.CORSmiddleware(handler.Login))
	mux.HandleFunc("/logout", handler.AuthMiddleware(handler.Logout))
	mux.HandleFunc("/passwords/create", api.CORSmiddleware(handler.AuthMiddleware(handler.CreatePasswordRequest)))
	mux.HandleFunc("/passwords/get", api.CORSmiddleware(handler.AuthMiddleware(handler.GetPasswordRequest)))
	mux.HandleFunc("/passwords/delete", api.CORSmiddleware(handler.AuthMiddleware(handler.DeletePasswordRequest)))
	mux.HandleFunc("/passwords/list", api.CORSmiddleware(handler.AuthMiddleware(handler.GetAllPasswordsRequest)))

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	slog.Info("Server is starting.", "address", serverAddr)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	server.Handler = mux

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
