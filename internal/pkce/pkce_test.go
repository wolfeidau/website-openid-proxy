package pkce

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MustNewVerifier(t *testing.T) {
	assert := require.New(t)

	verifier := MustNewVerifier(32)
	assert.Len(verifier, 43)
}

func Test_MustCodeChallengeS256(t *testing.T) {
	assert := require.New(t)

	challenge := MustCodeChallengeS256("GFRtrRUMZiEcFWlhW-3KxV4bBaQbj4T4pSCc_LjOuiE")
	assert.Len(challenge, 43)
	assert.Equal("KE-6gOh8H3HJ6cS28ZWoAWnHUFisTbK81AfSi6EP2gk", challenge)
}
