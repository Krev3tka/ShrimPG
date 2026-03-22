package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/Krev3tka/ShrimPG/internal/model"
)

type PasswordStorage interface {
	SavePassword(userID int, service, passwd, masterKey string) error
	GetPassword(serviceName, masterKey string) ([]byte, error)
	DeletePassword(service string) error
	VerifyMasterKey(ctx context.Context, userID int, masterKey string) (bool, error)
	GetAllPasswords(userID int, masterKey string) (model.Entry, error)
}

type Handler struct {
	storage       PasswordStorage
	sessions      map[string]string
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

func NewHandler(dbStorage PasswordStorage, userId int) *Handler {
	return &Handler{
		storage:       dbStorage,
		sessions:      make(map[string]string),
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
	if r.Body != nil {
		defer r.Body.Close()
	}

	var req struct {
		Masterkey string `json:"master_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ok, err := h.storage.VerifyMasterKey(r.Context(), h.currentUserID, req.Masterkey)
	if err != nil || !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := generateRandomToken()
	if err != nil {
		slog.Error("Failed to generate token.")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.mu.Lock()
	h.sessions[token] = req.Masterkey
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{"token": token}
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	delete(h.sessions, token)

	slog.Info("Logout successful", "token_prefix", token[:4])

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"logged out"}`))
}

func (h *Handler) CreatePasswordRequest(w http.ResponseWriter, r *http.Request) {
	masterKey, _ := r.Context().Value("masterKey").(string)

	var req SaveRequest
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.storage.SavePassword(h.currentUserID, req.Service, req.Password, masterKey); err != nil {
		http.Error(w, "Failed to save password"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	slog.Info("Password created successfully", "service", req.Service, "user_id", h.currentUserID)
}

func (h *Handler) GetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	slog.Info("GetPasswordRequest started", "method", r.Method)

	masterKey, ok := r.Context().Value("masterKey").(string)
	if !ok || masterKey == "" {
		http.Error(w, "Unauthorized: session missing key", http.StatusUnauthorized)
		return
	}

	slog.Info("Attempting to decode JSON body")
	var req ServiceRequest

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.Trim(req.Service, " ") == "" {
		slog.Error("Error: null service name", "user_id", h.currentUserID)
		http.Error(w, "Error: null service name", http.StatusNotFound)
		return
	}

	passwd, err := h.storage.GetPassword(req.Service, masterKey)
	if err != nil {
		http.Error(w, "Failed to get password: "+err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Password got successfully", "service", req.Service, "user_id", h.currentUserID)
	err = json.NewEncoder(w).Encode(PasswordResponse{
		Service:  req.Service,
		Password: string(passwd),
	})
	slog.Info("JSON decoded successfully", "service", req.Service)
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

func (h *Handler) GetAllPasswordsRequest(w http.ResponseWriter, r *http.Request) {
	masterKey, _ := r.Context().Value("masterKey").(string)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	passwords, err := h.storage.GetAllPasswords(h.currentUserID, masterKey)
	if err != nil {
		http.Error(w, "Failed to get passwords: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("All passwords retrieved", "count", len(passwords), "user_id", h.currentUserID)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(passwords); err != nil {
		slog.Error("Failed to encode response", "err", err)
		return
	}
}
