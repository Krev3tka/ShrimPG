package api

import (
	"context"
	"sync"

	"github.com/Krev3tka/ShrimPG/internal/model"
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
	storage  PasswordStorage
	sessions map[string]Session
	mu       sync.RWMutex
}

type Session struct {
	UserID int
	Key    string
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
