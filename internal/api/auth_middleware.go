// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("AUTH: Start checking token")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Warn("AUTH: Empty header")
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		data, err := h.rds.Get(r.Context(), "session:"+token).Bytes()
		if err != nil {
			http.Error(w, "Unauthorized: session expired", http.StatusUnauthorized)
			return
		}

		var sess Session
		if err := json.Unmarshal(data, &sess); err != nil {
			slog.Error("Failed to unmarshal session", "error", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		slog.Debug("Session data", "id", sess.UserID, "hasKey", sess.Key != "")

		ctx := context.WithValue(r.Context(), contextKey("userID"), sess.UserID)
		ctx = context.WithValue(ctx, contextKey("masterKey"), sess.Key)

		slog.Info("AUTH: Success, passing to handler", "user", sess.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
