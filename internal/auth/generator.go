// Copyright (C) 2026 krev3tka. Licensed under the GNU GPL v3.
package auth

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/Krev3tka/ShrimPG/internal/pkg/dictionary"
)

// GeneratePassphrase generates random passphrase consisting  of specified
// number of words joined by a '-' separator.
//
// It uses the Fisher-Yates shuffle algorithm and [crypto/rand] for
// cryptographically secure selection from the internal dictionary.
func GeneratePassphrase(length int) string {
	if length > len(dictionary.Words) {
		length = len(dictionary.Words)
	}
	shuffled := make([]string, len(dictionary.Words))
	copy(shuffled, dictionary.Words)

	for i := len(shuffled) - 1; i > 0; i-- {

		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue
		}
		j := nBig.Int64()

		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return strings.Join(shuffled[:length], "-")
}
