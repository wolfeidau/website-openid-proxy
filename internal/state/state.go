package state

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/labstack/echo/v4"
	"github.com/wolfeidau/aws-openid-proxy/internal/cookie"
)

const (
	StateKey    = "state_value"
	StateLength = 32
)

// Store provides a store for state values used during authentication
type Store interface {
	Init(c echo.Context) string
	Get(c echo.Context) (string, error)
}

var _ Store = &CookieStore{}

type CookieStore struct {
	cfg cookie.Config
}

func NewCookieStore(cfg cookie.Config) *CookieStore {
	return &CookieStore{cfg: cfg}
}

func (cs *CookieStore) Init(c echo.Context) string {
	state := mustRandomState()
	c.SetCookie(cookie.NewCookie(cs.cfg, state))
	// assign the state value to the context
	c.Set(StateKey, state)
	return state
}

func (cs *CookieStore) Get(c echo.Context) (string, error) {
	cookie, err := c.Cookie(cs.cfg.Name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

// Returns a base64 encoded random 32 byte string.
func mustRandomState() string {
	b := make([]byte, StateLength)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
