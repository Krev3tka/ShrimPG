// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package api

import (
	"errors"
	"net/http"

	"github.com/Krev3tka/ShrimPG/internal/crypto"
)

func getContextValues(r *http.Request) (int, error) {

	userID, ok := r.Context().Value(contextKey("userID")).(int)
	if !ok {
		return 0, errors.New("missing userID")
	}

	return userID, nil
}

// DecryptSessionKey decrypts the encrypted master key stored in session
func DecryptSessionKey(encryptedKey []byte, serverKey []byte) ([]byte, error) {
	if len(encryptedKey) == 0 {
		return nil, errors.New("empty encrypted key")
	}
	return crypto.Decrypt(encryptedKey, serverKey)
}
