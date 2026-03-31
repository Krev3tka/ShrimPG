// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package api

import (
	"context"

	"github.com/Krev3tka/ShrimPG/internal/model"
	"github.com/redis/go-redis/v9"
)

type PasswordStorage interface {
	SavePassword(userID int, service, passwd string, encryptionKey []byte) error
	GetPassword(userID int, serviceName string, encryptionKey []byte) ([]byte, error)
	DeletePassword(userID int, service string) error
	VerifyMasterKey(ctx context.Context, username string, masterKey string) (int, []byte, error)
	GetAllPasswords(userID int, encryptionKey []byte) (model.Entry, error)
	CreateUser(ctx context.Context, username string, masterKey string) (int, error)
}

type Handler struct {
	storage PasswordStorage
	rds     *redis.Client
}

type Session struct {
	UserID int    `json:"user_id"`
	Key    string `json:"key"`
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
