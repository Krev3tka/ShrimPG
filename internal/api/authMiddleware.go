package api

import (
	"net/http"
	"sync"
)

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		mu := sync.RWMutex{}
		mu.RLock()
		if !h.sessions[token] {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			mu.RUnlock()
			return
		}
		mu.RUnlock()
		next(w, r)
	}
}
