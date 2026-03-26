package api

import (
	"encoding/hex"
	"errors"
	"net/http"
)

func getContextValues(r *http.Request) (int, []byte, error) {
	keyHex, ok := r.Context().Value("masterKey").(string)
	if !ok {
		return 0, nil, errors.New("missing masterKey")
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		return 0, nil, errors.New("missing userID")
	}

	encryptionKey, err := hex.DecodeString(keyHex)
	if err != nil {
		return 0, nil, errors.New("invalid encryption key format")
	}

	return userID, encryptionKey, nil
}
