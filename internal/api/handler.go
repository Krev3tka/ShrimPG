// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package api

func NewHandler(dbStorage PasswordStorage) *Handler {
	return &Handler{
		storage:  dbStorage,
		sessions: make(map[string]Session),
	}
}
