package api

func NewHandler(dbStorage PasswordStorage) *Handler {
	return &Handler{
		storage:  dbStorage,
		sessions: make(map[string]Session),
	}
}
