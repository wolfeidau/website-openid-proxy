package session

import (
	"github.com/dghubble/sessions"
	"github.com/labstack/echo/v4"
	"github.com/wolfeidau/website-openid-proxy/internal/echosessions"
	"github.com/wolfeidau/website-openid-proxy/internal/flags"
	"github.com/wolfeidau/website-openid-proxy/internal/secrets"
)

// SetupMiddleware builds the session middleware after loading secrets
func SetupMiddleware(cfg *flags.API, secretCache *secrets.Cache) (echo.MiddlewareFunc, error) {

	sessionSecret, err := secretCache.GetValue(cfg.SessionSecretArn)
	if err != nil {
		return nil, err
	}

	// session middleware is available everwhere
	sessionMiddleware := echosessions.MiddlewareWithConfig(echosessions.Config{
		Store: sessions.NewCookieStore[string](
			sessions.DefaultCookieConfig,
			[]byte(sessionSecret),
			nil,
		),
	})

	return sessionMiddleware, nil
}
