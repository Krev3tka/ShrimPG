package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

type PasswordStorage interface {
	SavePassword(userID int, service, passwd, masterKey string) error
	GetPassword(serviceName, masterKey string) ([]byte, error)
	DeletePassword(service string) error
}

type Handler struct {
	storage       PasswordStorage
	masterKey     string
	sessions      map[string]bool
	mu            sync.RWMutex
	currentUserID int
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

func NewHandler(dbStorage PasswordStorage, key string, userId int) *Handler {
	return &Handler{
		storage:       dbStorage,
		masterKey:     key,
		sessions:      make(map[string]bool),
		currentUserID: userId,
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
	var req struct {
		Masterkey string `json:"master_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Masterkey != h.masterKey {
		http.Error(w, "Wrong master password", http.StatusForbidden)
		return
	}

	token, err := generateRandomToken()
	if err != nil {
		slog.Error("Failed to generate token.")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.mu.Lock()
	h.sessions[token] = true
	h.mu.Unlock()
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) CreatePasswordRequest(w http.ResponseWriter, r *http.Request) {
	var req SaveRequest

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.storage.SavePassword(h.currentUserID, req.Service, req.Password, h.masterKey); err != nil {
		http.Error(w, "Failed to save password"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	slog.Info("Password created successfully", "service", req.Service, "user_id", h.currentUserID)
}

func (h *Handler) GetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	var req ServiceRequest

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if strings.Trim(req.Service, " ") == "" {
		slog.Error("Error: null service name", "user_id", h.currentUserID)
		http.Error(w, "Error: null service name", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	passwd, err := h.storage.GetPassword(req.Service, h.masterKey)
	if err != nil {
		http.Error(w, "Failed to get password: "+err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Password got successfully", "service", req.Service, "user_id", h.currentUserID)
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
		slog.Info("Password already deleted or it didn't exist", "service", "user_id", h.currentUserID)
	}

	err = h.storage.DeletePassword(req.Service)
	if err != nil {
		slog.Error("Failed to delete", "service", req.Service, "err", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Password deleted successfully", "user_id", h.currentUserID)
}
