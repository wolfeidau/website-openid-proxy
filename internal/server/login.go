package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/s3website-openid-proxy/internal/echosessions"
	"github.com/wolfeidau/s3website-openid-proxy/internal/pkce"
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

// UserInfo user info returned by user info route
type UserInfo struct {
	Sub   string `json:"sub,omitempty"`
	Email string `json:"email,omitempty"`
}

func userInfoFromValues(val map[string]interface{}) (*UserInfo, error) {
	sub, ok := val["sub"].(string)
	if !ok {
		return nil, errors.New("failed to read sub")
	}
	email, ok := val["email"].(string)
	if !ok {
		return nil, errors.New("failed to read email")
	}

	return &UserInfo{
		Sub:   sub,
		Email: email,
	}, nil
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
	verifier := pkce.MustNewVerifier(32)

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
	authSess.Values["verifier"] = verifier

	err = authSess.Save(c.Response())
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to save session")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	redirectURL := l.oauthConfig().AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", pkce.MustCodeChallengeS256(verifier)),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	// send the caller off to their login server
	return c.Redirect(http.StatusFound, redirectURL)
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

	verifier, ok := authSess.Values["verifier"].(string)
	if !ok {
		log.Ctx(ctx).Error().Msg("missing verifier attribute from session")

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

	tokens, err := l.oauthConfig().Exchange(ctx, cb.Code, oauth2.SetAuthURLParam("code_verifier", verifier))
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

	ctx := c.Request().Context()

	session, err := echosessions.Get(loggedInCookieName, c)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get session")

		// TODO: Need an error page
		return c.String(http.StatusUnauthorized, "failed to process request")
	}

	info, err := userInfoFromValues(session.Values)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to read user info from session")

		// TODO: Need an error page
		return c.String(http.StatusInternalServerError, "failed to process request")
	}

	return c.JSON(http.StatusOK, info)
}

// Logout logout http handler
func (l *Auth) Logout(c echo.Context) error {

	ctx := c.Request().Context()

	err := echosessions.Destroy(loggedInCookieName, c)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to destroy session")

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
