// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package api

import (
	"github.com/redis/go-redis/v9"
)

func NewHandler(dbStorage PasswordStorage, redisClient *redis.Client) *Handler {
	return &Handler{
		storage: dbStorage,
		rds:     redisClient,
	}
}
