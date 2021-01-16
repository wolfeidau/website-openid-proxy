package server

import (
	"context"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/s3website-openid-proxy/internal/echosessions"
	"golang.org/x/oauth2"
)

const (
	authCookieName   = "proxy_auth_session"
	authCookieExpiry = 5 * 60 // 5 minutes

	loggedInCookieName   = "proxy_login_session"
	loggedInCookieExpiry = 8 * 60 * 60 // 8 hours

	stateLength = 32
)

// Callback callback info
type Callback struct {
	Code  string `query:"code"`
	State string `query:"state"`
}

// Auth authentication related handlers
type Auth struct {
	authConfig *AuthConfig
	provider   *oidc.Provider
}

// NewAuth new auth server http handlers
func NewAuth(ac *AuthConfig) (*Auth, error) {

	if err := ac.Valid(); err != nil {
		return nil, err
	}

	provider, err := ac.ProviderFunc(context.Background(), ac.Issuer)
	if err != nil {
		return nil, err
	}

	return &Auth{authConfig: ac, provider: provider}, nil
}

// Login login http handler
func (l *Auth) Login(c echo.Context) error {

	ctx := c.Request().Context()

	state := MustRandomState(stateLength)

	authSess, err := echosessions.New(authCookieName, c)
	if err != nil {
		log.Ctx(ctx).Error().Stack().Err(err).Msg("failed to create new session")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	// override the default cookie settings
	authSess.Config.MaxAge = authCookieExpiry
	authSess.Config.Secure = true
	authSess.Values["state"] = state

	err = authSess.Save(c.Response())
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to save session")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	// send the caller off to their login server
	return c.Redirect(http.StatusFound, l.oauthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline))
}

// Callback callback http handler
func (l *Auth) Callback(c echo.Context) error {

	ctx := c.Request().Context()

	cb := new(Callback)

	if err := c.Bind(cb); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to Bind session form")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	authSess, err := echosessions.Get(authCookieName, c)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get auth session")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	state, ok := authSess.Values["state"].(string)
	if !ok {
		log.Ctx(ctx).Error().Msg("missing state attribute from session")

		// TODO: Need an error page
		return c.String(http.StatusBadRequest, "failed to process request")
	}

	// clean up the completed auth session
	defer authSess.Destroy(c.Response())

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

	loginSess, err := echosessions.New(loggedInCookieName, c)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get new session")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	// override the default cookie settings
	loginSess.Config.Secure = true
	loginSess.Config.MaxAge = loggedInCookieExpiry
	loginSess.Values["email"] = userInfo.Email
	loginSess.Values["sub"] = userInfo.Subject

	err = loginSess.Save(c.Response())
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to save session")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	return c.Redirect(http.StatusFound, "/")
}

// UserInfo user info http handler
func (l *Auth) UserInfo(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

// Logout logout http handler
func (l *Auth) Logout(c echo.Context) error {

	ctx := c.Request().Context()

	err := echosessions.Destroy(loggedInCookieName, c)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get new session")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	return c.NoContent(http.StatusOK)
}

// RegisterRoutes register the login related auth routes
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
