// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package api

import (
	"github.com/Krev3tka/ShrimPG/internal/db"
	"github.com/redis/go-redis/v9"
)

func NewHandler(dbStorage *db.DBStorage, redisClient *redis.Client, serverKey []byte) *Handler {
	return &Handler{
		storage:   dbStorage,
		rds:       redisClient,
		serverKey: serverKey,
	}
}
