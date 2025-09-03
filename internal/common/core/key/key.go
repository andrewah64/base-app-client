package key

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
)

func Key(length int) (string, error) {
	bytes := make([]byte,length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes), nil
}

func Hash (key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}
