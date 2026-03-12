package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/Krev3tka/ShrimPG/internal/crypto"
	"github.com/Krev3tka/ShrimPG/internal/db"
)

type Handler struct {
	storage   *db.DBStorage
	masterKey string
	sessions  map[string]bool
	mu        sync.RWMutex
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

func NewHandler(dbStorage *db.DBStorage, key string) *Handler {
	return &Handler{
		storage:   dbStorage,
		masterKey: key,
		sessions:  make(map[string]bool),
	}
}

func generateRandomToken() (string, error) {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	token, err := generateRandomToken()
	if err != nil {
		slog.Error("Failed to generate token.")
		return
	}
	h.mu.Lock()
	h.sessions[token] = true
	h.mu.Unlock()
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) CreatePasswordRequest(w http.ResponseWriter, r *http.Request) {
	p := &crypto.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  12,
		KeyLength:   16,
	}
	var req SaveRequest

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.storage.SavePassword(1, req.Service, req.Password, h.masterKey, p); err != nil {
		http.Error(w, "Failed to save password"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	slog.Info("Password created successfully", "service", req.Service, "user_id", 1)
}

func (h *Handler) GetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	p := &crypto.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  12,
		KeyLength:   16,
	}
	var req ServiceRequest
	salt, err := crypto.GenerateRandomBytes(p.SaltLength)
	if err != nil {
		http.Error(w, "Failed to generate salt", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if strings.Trim(req.Service, " ") == "" {
		slog.Error("Error: null service name", "user_id", 1)
		http.Error(w, "Error: null service name", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	passwd, err := h.storage.GetPassword(req.Service, h.masterKey, salt, p)
	if err != nil {
		http.Error(w, "Failed to get password: "+err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Password got successfully", "service", req.Service, "user_id", 1)
	err = json.NewEncoder(w).Encode(PasswordResponse{
		Service:  req.Service,
		Password: string(passwd),
	})
	if err != nil {
		http.Error(w, "Failed to write JSON"+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) DeletePasswordRequest(w http.ResponseWriter, r *http.Request) {
	var req ServiceRequest

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.Trim(req.Service, " ") == "" {
		slog.Info("Password already deleted or it didn't exist", "service", "user_id", 1)
	}

	err = h.storage.DeletePassword(req.Service)
	if err != nil {
		slog.Error("Failed to delete", "service", req.Service, "err", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Password deleted successfully", "user_id", 1)
}
