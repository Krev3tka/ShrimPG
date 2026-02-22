package storage

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func GenerateKoolPassword(length int) string {

	if length > len(words) {
		length = len(words)
	}
	shuffled := make([]string, len(words))
	copy(shuffled, words)

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
