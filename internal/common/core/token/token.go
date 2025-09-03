package token

import (
	"crypto/rand"
	"encoding/base64"
)

func Token(length int) (string, error) {
	bytes := make([]byte,length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
