package server

import (
	"context"

	"github.com/coreos/go-oidc"
)

type ProviderFunc = func(ctx context.Context, issuer string) (*oidc.Provider, error)
