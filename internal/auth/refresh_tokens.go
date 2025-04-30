package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	randomData, err := rand.Read(make([]byte, 32))
	if err != nil {
		return "", err
	}

	encodedData := hex.EncodeToString(make([]byte, randomData))

	return encodedData, nil
}
