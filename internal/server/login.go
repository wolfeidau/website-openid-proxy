package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/dghubble/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/aws-openid-proxy/internal/state"
	"golang.org/x/oauth2"
)

type ProviderFunc = func(ctx context.Context, issuer string) (*oidc.Provider, error)

// AuthConfig authentication configuration for to enable openid login
type AuthConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	StateStore   state.Store
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

	if ac.StateStore == nil {
		return errors.New("empty StateStore")
	}

	if ac.ProviderFunc == nil {
		ac.ProviderFunc = oidc.NewProvider
	}

	return nil
}

// Auth authentication related handlers
type Auth struct {
	authConfig   *AuthConfig
	sessionStore sessions.Store
	provider     *oidc.Provider
}

func NewAuth(ac *AuthConfig, sessionStore sessions.Store) (*Auth, error) {

	if err := ac.Valid(); err != nil {
		return nil, err
	}

	provider, err := ac.ProviderFunc(context.Background(), ac.Issuer)
	if err != nil {
		return nil, err
	}

	return &Auth{authConfig: ac, sessionStore: sessionStore, provider: provider}, nil
}

func (l *Auth) Login(c echo.Context) error {

	state := l.authConfig.StateStore.Init(c)

	// send the caller off to their login server
	return c.Redirect(http.StatusFound, l.oauthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline))
}

func (l *Auth) Callback(c echo.Context) error {

	ctx := c.Request().Context()

	cb := new(Callback)

	if err := c.Bind(cb); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to Bind session form")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	state, err := l.authConfig.StateStore.Get(c)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get state")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	if cb.State == "" && state != cb.State {
		log.Ctx(ctx).Error().Err(err).Msg("failed to validate state")

		// TODO: Need an error page
		return c.String(http.StatusBadRequest, "failed to process request")
	}

	tokens, err := l.oauthConfig().Exchange(ctx, cb.Code)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to exchange tokens")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	userInfo, err := l.provider.UserInfo(ctx, oauth2.StaticTokenSource(tokens))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get userinfo")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	sess := l.sessionStore.New("proxy")

	sess.Values["email"] = userInfo.Email
	sess.Values["sub"] = userInfo.Subject

	err = sess.Save(c.Response())
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to save session")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	return c.Redirect(http.StatusFound, "/")
}

func (l *Auth) UserInfo(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (l *Auth) Logout(c echo.Context) error {

	l.sessionStore.Destroy(c.Response(), "proxy")

	return c.NoContent(http.StatusOK)
}

func (l *Auth) RegisterRoutes(r interface {
	GET(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
}) {
	r.GET("/login", l.Login)
	r.GET("/callback", l.Callback)
	r.GET("/userinfo", l.UserInfo)
	r.GET("/logout", l.Logout)
}

func (l *Auth) oauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     l.authConfig.ClientID,
		ClientSecret: l.authConfig.ClientSecret,
		Endpoint:     l.provider.Endpoint(),
		RedirectURL:  l.authConfig.RedirectURL,
		Scopes:       []string{"email", "openid"},
	}
}

// LoginSkipper used to avoid running middleware for login requests
func LoginSkipper(prefix string) middleware.Skipper {
	return func(c echo.Context) bool {
		return strings.HasPrefix(c.Path(), prefix)
	}
}

// Callback callback info
type Callback struct {
	Code  string `query:"code"`
	State string `query:"state"`
}
