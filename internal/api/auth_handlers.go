// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package api

import (
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Krev3tka/ShrimPG/internal/auth"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req struct {
		Username  string `json:"username"`
		Masterkey string `json:"masterKey"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, derivedKey, err := h.storage.VerifyMasterKey(r.Context(), req.Username, req.Masterkey)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateRandomToken()
	if err != nil || token == "" {
		http.Error(w, "Failed to generate token", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	h.sessions[token] = Session{
		userID,
		hex.EncodeToString(derivedKey),
	}
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{"token": token}
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	h.mu.Lock()
	delete(h.sessions, token)
	h.mu.Unlock()

	slog.Info("Logout attempt", "token_len", len(token))
	if len(token) > 4 {
		slog.Info("Logout successful", "token_prefix", token[:4])
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"logged out"}`))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req struct {
		Username  string `json:"username"`
		Masterkey string `json:"masterKey"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Username) <= 6 {
		slog.Warn("Registration failed: username too short", "username", req.Username)
	}

	if _, err := h.storage.CreateUser(r.Context(), req.Username, req.Masterkey); err != nil {
		http.Error(w, "Failed to create user", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
