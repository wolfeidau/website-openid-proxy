package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/go-oidc"
	"github.com/dghubble/sessions"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/website-openid-proxy/internal/echosessions"
	"github.com/wolfeidau/website-openid-proxy/internal/flags"
	"github.com/wolfeidau/website-openid-proxy/internal/logger"
	"github.com/wolfeidau/website-openid-proxy/mocks"
)

var mockProviderFunc = func(ctx context.Context, issuer string) (*oidc.Provider, error) {
	return &oidc.Provider{}, nil
}

func TestLogin(t *testing.T) {

	assert := require.New(t)

	auth, err := NewAuth(newConfig(), mockProviderFunc)
	assert.NoError(err)

	e := echo.New()

	sessionMiddleware := echosessions.MiddlewareWithConfig(echosessions.Config{
		Store: sessions.NewCookieStore[string](sessions.DefaultCookieConfig, []byte("test"), nil),
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
	assert.Contains(rec.Result().Header.Get(echo.HeaderLocation), "redirect_uri=http%3A%2F%2Flocalhost%2Fcallback&response_type=code")
}

func TestUserInfo(t *testing.T) {
	assert := require.New(t)

	auth, err := NewAuth(newConfig(), mockProviderFunc)
	assert.NoError(err)

	e := echo.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionStore := mocks.NewMockStore(ctrl)

	sess := sessions.NewSession[string](sessionStore, loggedInCookieName)

	sess.Set("sub", "abc123")
	sess.Set("email", "mark@wolfe.id.au")

	sessionStore.EXPECT().Get(gomock.Any(), loggedInCookieName).Return(sess, nil)

	sessionMiddleware := echosessions.MiddlewareWithConfig(echosessions.Config{
		Store: sessionStore,
	})

	req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
	req = req.WithContext(logger.NewLoggerWithContext(context.TODO()))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// wrap the call in the required session middleware
	h := sessionMiddleware(func(c echo.Context) error {
		return auth.UserInfo(c)
	})

	// call the handler func with the context
	err = h(c)

	assert.NoError(err)
	assert.Equal(http.StatusOK, rec.Result().StatusCode)
	assert.JSONEq(`{"sub":"abc123","email":"mark@wolfe.id.au"}`, rec.Body.String())
}

func TestUserInfo_StatusUnauthorized(t *testing.T) {
	assert := require.New(t)

	auth, err := NewAuth(newConfig(), mockProviderFunc)
	assert.NoError(err)

	e := echo.New()

	sessionMiddleware := echosessions.MiddlewareWithConfig(echosessions.Config{
		Store: sessions.NewCookieStore[string](sessions.DefaultCookieConfig, []byte("test"), nil),
	})

	req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
	req = req.WithContext(logger.NewLoggerWithContext(context.TODO()))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// wrap the call in the required session middleware
	h := sessionMiddleware(func(c echo.Context) error {
		return auth.UserInfo(c)
	})

	// call the handler func with the context
	err = h(c)

	assert.NoError(err)
	assert.Equal(http.StatusUnauthorized, rec.Result().StatusCode)
}

func newConfig() *flags.API {
	return &flags.API{
		Issuer:       "http://localhost",
		ClientID:     "abc123",
		ClientSecret: "cde456",
		RedirectURL:  "http://localhost/callback",
	}
}
