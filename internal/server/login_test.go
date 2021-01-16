package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/go-oidc"
	"github.com/dghubble/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/s3website-openid-proxy/internal/echosessions"
	"github.com/wolfeidau/s3website-openid-proxy/internal/logger"
)

func TestLogin(t *testing.T) {

	assert := require.New(t)

	auth, err := NewAuth(&AuthConfig{
		Issuer:       "http://localhost",
		ClientID:     "abc123",
		ClientSecret: "cde456",
		RedirectURL:  "http://localhost/callback",
		ProviderFunc: func(ctx context.Context, issuer string) (*oidc.Provider, error) {
			return &oidc.Provider{}, nil
		},
	})
	assert.NoError(err)

	e := echo.New()

	sessionMiddleware := echosessions.MiddlewareWithConfig(echosessions.Config{
		Store: sessions.NewCookieStore([]byte("test"), nil),
	})

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req = req.WithContext(logger.NewLoggerWithContext(context.TODO()))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// wrap the call in the required session middleware
	h := sessionMiddleware(func(c echo.Context) error {
		return auth.Login(c)
	})

	// call the handler func with the context
	err = h(c)

	assert.NoError(err)
	assert.Equal(302, rec.Result().StatusCode)
	assert.Contains(rec.Result().Header.Get(echo.HeaderLocation), "client_id=abc123&redirect_uri=http%3A%2F%2Flocalhost%2Fcallback&response_type=code")
}
