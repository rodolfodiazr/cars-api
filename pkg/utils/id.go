package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID creates a random 16-character hexadecimal ID.
func GenerateID() (string, error) {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
