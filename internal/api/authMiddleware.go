package api

import (
	"net/http"
)

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//token := r.Header.Get("Authorization")
		//h.mu.RLock()
		//if token == "" {
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	h.mu.RUnlock()
		//	return
		//}
		//valid := h.sessions[token]
		//h.mu.RUnlock()
		//
		//if !valid {
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	return
		//}
		next(w, r)
	}
}
