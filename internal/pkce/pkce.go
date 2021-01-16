package pkce

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// MustNewVerifier generate and encode a new verifier
func MustNewVerifier(verifierLength int) string {
	b := make([]byte, verifierLength)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return encode(b)
}

// MustCodeChallengeS256 using an already encoded verifier generate the challenge
func MustCodeChallengeS256(verifier string) string {
	h := sha256.New()

	_, err := h.Write([]byte(verifier))
	if err != nil {
		panic(err)
	}

	return encode(h.Sum(nil))
}

func encode(msg []byte) string {
	return base64.RawURLEncoding.EncodeToString(msg)
}
