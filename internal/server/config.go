package server

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc"
)

type ProviderFunc = func(ctx context.Context, issuer string) (*oidc.Provider, error)

// AuthConfig authentication configuration for to enable openid login
type AuthConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	ProviderFunc ProviderFunc
}

// Valid validate our configuration
func (ac *AuthConfig) Valid() error {
	if ac.Issuer == "" {
		return errors.New("empty Issuer")
	}
	if ac.ClientID == "" {
		return errors.New("empty ClientID")
	}
	if ac.ClientSecret == "" {
		return errors.New("empty ClientSecret")
	}

	if ac.RedirectURL == "" {
		return errors.New("empty RedirectURL")
	}

	if ac.ProviderFunc == nil {
		ac.ProviderFunc = oidc.NewProvider
	}

	return nil
}
