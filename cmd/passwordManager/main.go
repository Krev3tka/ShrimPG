// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package main

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Krev3tka/ShrimPG-backend/internal/api"
	"github.com/Krev3tka/ShrimPG-backend/internal/crypto"
	"github.com/Krev3tka/ShrimPG-backend/internal/db"
	"github.com/Krev3tka/ShrimPG-backend/internal/storage"
	"github.com/joho/godotenv"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system env")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Generate server encryption key (persistent across server lifetime)
	encodedKey := os.Getenv("SERVER_SECRET_KEY")
	var serverKey []byte

	if encodedKey != "" {
		decoded, err := base64.StdEncoding.DecodeString(encodedKey)
		if err != nil || len(decoded) != 32 {
			slog.Error("Invalid SERVER_SECRET_KEY in .env (must be 32 bytes base64)", "error", err)
			return
		}
		serverKey = decoded
		slog.Info("Using persistent server key from .env")
	} else {
		slog.Warn("SERVER_SECRET_KEY not found in .env, generating temporary key")
		serverKey, _ = crypto.GenerateRandomBytes(32)
	}

	// __Redis Configuration__

	// Get Redis address from environment variables
	rdsAddr := getEnv("REDISDB_ADDRESS", "127.0.0.1:6379")
	slog.Info("Checking Redis address", "addr", rdsAddr)

	// Initialize Redis configuration
	cfg := storage.Config{
		Addr:        rdsAddr,
		Password:    getEnv("REDISDB_PASSWORD", ""),
		User:        "",
		DB:          0,
		MaxRetries:  3,
		DialTimeout: 5 * time.Second,
		Timeout:     10 * time.Second,
	}

	// Initialize Redis client
	redisDb, err := storage.NewClient(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	// Create database connection pool
	dbPool, err := db.Connect()
	if err != nil {
		slog.Error("Database connection failed", "details", err)
		return
	}

	// __Server Address Logic__

	// Get HTTP server address from environment variables
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "127.0.0.1:8080"
	}

	u, err := url.Parse(serverAddr)
	if err == nil {
		slog.Info("database connection established", "host", u.Host)
	}

	// Initialize vault storage
	vault := db.NewDBStorage(dbPool)

	// Database readiness check. Required if the backend starts faster than the database container
	for i := 0; i < 5; i++ {
		pingCtx, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
		err = dbPool.Ping(pingCtx)
		cancelPing()
		if err == nil {
			break
		}
		slog.Info("Waiting for database...", "attempt", i+1)
		time.Sleep(2 * time.Second)
	}

	// Initialize SQL migrations
	slog.Info("running migrations...")
	if err := db.InitMigrations(dbPool, "./migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		return
	}

	// __Handlers__

	// Create API handler
	handler := api.NewHandler(vault, redisDb, serverKey)

	// Initialize endpoints and routing
	mux := http.NewServeMux()
	mux.HandleFunc("/register", api.CORSmiddleware(handler.Register))
	mux.HandleFunc("/login", api.CORSmiddleware(handler.RateLimit(handler.Login)))
	mux.HandleFunc("/logout", handler.AuthMiddleware(handler.Logout))
	mux.HandleFunc("/passwords/create", api.CORSmiddleware(handler.AuthMiddleware(handler.CreatePasswordRequest)))
	mux.HandleFunc("/passwords/get", api.CORSmiddleware(handler.AuthMiddleware(handler.GetPasswordRequest)))
	mux.HandleFunc("/passwords/delete", api.CORSmiddleware(handler.AuthMiddleware(handler.DeletePasswordRequest)))
	mux.HandleFunc("/passwords/list", api.CORSmiddleware(handler.AuthMiddleware(handler.GetAllPasswordsRequest)))

	// Setup for graceful shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	slog.Info("Server is starting.", "address", serverAddr)

	// __Server__

	// Initialize HTTP server settings
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Goroutine for server execution. Handles potential crashes
	go func() {
		if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil && err != http.ErrServerClosed {
			slog.Error("Server crashed unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown handling

	<-exit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	dbPool.Close()
	_ = redisDb.Close()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		slog.Error("Error: failed to shutdown server correctly", "error", err)
		return
	}
	slog.Info("Server stopped gracefully")
}
