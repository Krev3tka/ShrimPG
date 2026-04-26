// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package api

import (
	"context"

	"github.com/Krev3tka/ShrimPG/internal/model"
	"github.com/redis/go-redis/v9"
)

type PasswordStorage interface {
	SavePassword(userID int, service, encryptedData []byte) error
	GetPassword(userID int, serviceName string) ([]byte, error)
	DeletePassword(userID int, service string) error
	VerifyAuthHash(ctx context.Context, username string, masterKey string) (int, error)
	GetAllPasswords(userID int) (model.Entry, error)
	CreateUser(ctx context.Context, username string, masterKey string) (int, error)
}

type Handler struct {
	storage   PasswordStorage
	rds       *redis.Client
	serverKey []byte
}

type Session struct {
	UserID       int    `json:"user_id"`
	EncryptedKey []byte `json:"encrypted_key"`
}

type SaveRequest struct {
	Service  string `json:"service"`
	Password string `json:"password"`
}

type ServiceRequest struct {
	Service string `json:"service"`
}

type PasswordResponse struct {
	Service  string `json:"service"`
	Password string `json:"password"`
}

type contextKey string

type AuthRequest struct {
	Username string `json:"username"`
	AuthHash string `json:"authHash"` // Результат Argon2/PBKDF2 на клиенте
}
