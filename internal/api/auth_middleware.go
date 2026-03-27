// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
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
		session, exists := h.sessions[token]
		h.mu.RUnlock()

		if !exists {
			http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKey("masterKey"), session.Key)
		ctx = context.WithValue(ctx, contextKey("userID"), session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
