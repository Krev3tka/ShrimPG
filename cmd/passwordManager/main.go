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
	"github.com/Krev3tka/ShrimPG/internal/crypto/rpc"
	"github.com/Krev3tka/ShrimPG/internal/db"
	"github.com/Krev3tka/ShrimPG/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// __Конфигурация редиса__

	// Получение редис-БД адрФеса
	rdsAddr := os.Getenv("REDISDB_ADDRESS")

	// Инициализация редис-конфигурации
	cfg := storage.Config{
		Addr:        rdsAddr,
		Password:    "",
		User:        "",
		DB:          0,
		MaxRetries:  3,
		DialTimeout: 1 * time.Second,
		Timeout:     10 * time.Second,
	}

	// Инициализация клиента редис БД
	redisDb, err := storage.NewClient(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	// Создание пула соединений для БД
	dbPool, err := db.Connect()
	if err != nil {
		slog.Error("Database connection failed", "details", err)
		return
	}

	// __Логика порта__

	// Получение адреса http-сервера
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "0.0.0.0:8080"
	}

	u, err := url.Parse(serverAddr)
	if err == nil {
		slog.Info("database connection established", "host", u.Host)
	}

	// Инициализация хранилища
	vault := db.NewDBStorage(dbPool)
	if cryptoAddr := os.Getenv("CRYPTO_SERVICE_ADDR"); cryptoAddr != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn, err := grpc.DialContext(ctx, cryptoAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		cancel()
		if err != nil {
			slog.Warn("crypto service unavailable; falling back to local crypto", "error", err)
		} else {
			vault = db.NewDBStorageWithCrypto(dbPool, db.NewCryptoEngineFromConn(rpc.NewCryptoServiceClient(conn)))
			defer conn.Close()
			slog.Info("crypto service connected", "address", cryptoAddr)
		}
	}

	// Проверка БД на готовность. Эта проверка нужна на случай, когда Go-backend запускается быстрее БД
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

	// Инициализация SQL-миграций
	slog.Info("running migrations...")
	if err := db.InitMigrations(dbPool, "./migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		return
	}

	// __Хендлеры__

	// Создание хендлера
	handler := api.NewHandler(vault, redisDb)

	// Инициализация эндпоинтов
	mux := http.NewServeMux()
	mux.HandleFunc("/register", api.CORSmiddleware(handler.Register))
	mux.HandleFunc("/login", api.CORSmiddleware(handler.Login))
	mux.HandleFunc("/logout", handler.AuthMiddleware(handler.Logout))
	mux.HandleFunc("/passwords/create", api.CORSmiddleware(handler.AuthMiddleware(handler.CreatePasswordRequest)))
	mux.HandleFunc("/passwords/get", api.CORSmiddleware(handler.AuthMiddleware(handler.GetPasswordRequest)))
	mux.HandleFunc("/passwords/delete", api.CORSmiddleware(handler.AuthMiddleware(handler.DeletePasswordRequest)))
	mux.HandleFunc("/passwords/list", api.CORSmiddleware(handler.AuthMiddleware(handler.GetAllPasswordsRequest)))

	// Заготовка для graceful shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	slog.Info("Server is starting.", "address", serverAddr)

	// __Сервер__

	// Инициализация сервера
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	server.Handler = mux

	// Горутина для работы сервера. Обработка ошибки и потенциального краша
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server crashed unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown

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
