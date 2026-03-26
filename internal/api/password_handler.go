package api

import (
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handler) CreatePasswordRequest(w http.ResponseWriter, r *http.Request) {
	userID, encryptionKey, err := getContextValues(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

	if err := h.storage.SavePassword(userID, req.Service, req.Password, encryptionKey); err != nil {
		http.Error(w, "Failed to save password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	slog.Info("Password created successfully", "service", req.Service, "user_id", userID)
}

func (h *Handler) GetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	userID, encryptionKey, err := getContextValues(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req ServiceRequest

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.Trim(req.Service, " ") == "" {
		slog.Error("Error: null service name", "user_id", userID)
		http.Error(w, "Error: null service name", http.StatusBadRequest)
		return
	}

	if len(req.Service) > 40 {
		slog.Error("Error: too long service name", "user_id", userID)
		http.Error(w, "Error: too long service name", http.StatusBadRequest)
		return
	}

	passwd, err := h.storage.GetPassword(userID, req.Service, encryptionKey)
	if err != nil {
		http.Error(w, "Failed to get password: "+err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Password got successfully", "service", req.Service, "user_id", userID)
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

	userID, _, err := getContextValues(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.Trim(req.Service, " ") == "" {
		slog.Info("Password already deleted or it didn't exist", "service", "user_id", userID)
	}

	err = h.storage.DeletePassword(userID, req.Service)
	if err != nil {
		slog.Error("Failed to delete", "service", req.Service, "err", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Password deleted successfully", "user_id", userID)
}

func (h *Handler) GetAllPasswordsRequest(w http.ResponseWriter, r *http.Request) {
	keyHex, ok := r.Context().Value("masterKey").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	encryptionKey, err := hex.DecodeString(keyHex)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	passwords, err := h.storage.GetAllPasswords(userID, encryptionKey)
	if err != nil {
		http.Error(w, "Failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("All passwords retrieved", "count", len(passwords), "user_id", userID)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(passwords); err != nil {
		slog.Error("Failed to encode response", "err", err)
		return
	}
}
