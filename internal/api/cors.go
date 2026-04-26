package api

import (
	"net/http"
)

// CORSmiddleware manages Cross-Origin Resource Sharing for trusted domains.
func CORSmiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Whitelist of trusted origins
		allowedOrigins := map[string]bool{
			"http://localhost:5173": true,
			"http://localhost:3000": true,
			"https://shrimpg.app":   true,
		}

		// Get the Origin header from the incoming request
		origin := r.Header.Get("Origin")

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}
