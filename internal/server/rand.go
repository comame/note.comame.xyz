package server

import (
	"crypto/rand"
	"encoding/base64"
)

func randomString(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b64 := base64.RawURLEncoding.EncodeToString(b)
	return b64[:size], nil
}
