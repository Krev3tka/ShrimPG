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

	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "0.0.0.0:8080"
	}

	slog.Info("database connection established", "address", os.Getenv("CONN_STR"))

	vault := db.NewDBStorage(dbPool)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		err = dbPool.Ping(context.Background())
		if err == nil {
			break
		}
		slog.Info("Waiting for database...", "attempt", i+1)
		time.Sleep(2 * time.Second)
	}

	if err := vault.Ping(ctx); err != nil {
		slog.Error("Database is unreachable", "error", err)
		return
	}

	slog.Info("running migrations...")
	if err := db.InitMigrations(dbPool, "./migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		return
	}

	handler := api.NewHandler(vault)

	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/logout", handler.AuthMiddleware(handler.Logout))
	http.HandleFunc("/passwords/create", handler.AuthMiddleware(handler.CreatePasswordRequest))
	http.HandleFunc("/passwords/get", handler.AuthMiddleware(handler.GetPasswordRequest))
	http.HandleFunc("/passwords/delete", handler.AuthMiddleware(handler.DeletePasswordRequest))
	http.HandleFunc("/passwords/list", handler.AuthMiddleware(handler.GetAllPasswordsRequest))

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
