package api

import (
	"context"
	"net/http"
	"strings"
)

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		h.mu.RLock()
		masterKey, exists := h.sessions[token]
		defer h.mu.RUnlock()

		if !exists || masterKey == "" {
			http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "masterKey", masterKey)
		next(w, r.WithContext(ctx))
	}
}
