// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package api

import (
	"crypto/sha256"
	"errors"
	"net/http"
)

func getContextValues(r *http.Request) (int, []byte, error) {
	keyHex, ok := r.Context().Value(contextKey("masterKey")).(string)
	if !ok {
		return 0, nil, errors.New("missing masterKey")
	}

	userID, ok := r.Context().Value(contextKey("userID")).(int)
	if !ok {
		return 0, nil, errors.New("missing userID")
	}

	hash := sha256.Sum256([]byte(keyHex))
	encryptionKey := hash[:]

	return userID, encryptionKey, nil
}
