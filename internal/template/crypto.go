package template

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

func sha256sum(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func sha1sum(input string) string {
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// uuidv4 provides a safe and secure UUID v4 implementation
func uuidv4() string {
	return uuid.New().String()
}
