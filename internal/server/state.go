package server

import (
	"crypto/rand"
	"encoding/base64"
)

// MustRandomState returns a base64 encoded random 32 byte string.
func MustRandomState(len int) string {
	b := make([]byte, len)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
