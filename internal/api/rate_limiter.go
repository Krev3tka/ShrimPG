package api

import (
	"log/slog"
	"net"
	"net/http"
	"time"
)

func (h *Handler) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		key := "limiter:login:" + ip

		count, err := h.rds.Incr(r.Context(), key).Result()
		if err != nil {
			slog.Error("Redis limiter error", "err", err)
			next(w, r)
			return
		}

		if count == 1 {
			h.rds.Expire(r.Context(), key, time.Minute)
		}

		if count > 5 {
			slog.Warn("Rate limit exceeded", "ip", ip)
			http.Error(w, "Too many attempts. Try again in a minute.", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
