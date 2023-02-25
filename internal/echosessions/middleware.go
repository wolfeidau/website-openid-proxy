package echosessions

import (
	"fmt"

	"github.com/dghubble/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// Config defines the config for Session middleware.
	Config struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Session store.
		// Required.
		Store sessions.Store[string]
	}
)

const (
	key = "_session_store"
)

var (
	// DefaultConfig is the default Session middleware config.
	DefaultConfig = Config{
		Skipper: middleware.DefaultSkipper,
	}
)

// Get returns a named session.
func Get(name string, c echo.Context) (*sessions.Session[string], error) {
	s := c.Get(key)
	if s == nil {
		return nil, fmt.Errorf("%q session not found", name)
	}
	store := s.(sessions.Store[string])
	return store.Get(c.Request(), name)
}

// New returns a new session
func New(name string, c echo.Context) (*sessions.Session[string], error) {
	s := c.Get(key)
	if s == nil {
		return nil, fmt.Errorf("%q session not found", name)
	}
	store := s.(sessions.Store[string])

	return store.New(name), nil
}

// Destroy destroy existing session
func Destroy(name string, c echo.Context) error {
	s := c.Get(key)
	if s == nil {
		return fmt.Errorf("%q session not found", name)
	}
	store := s.(sessions.Store[string])

	store.Destroy(c.Response(), name)

	return nil
}

// Middleware returns a Session middleware.
func Middleware(store sessions.Store[string]) echo.MiddlewareFunc {
	c := DefaultConfig
	c.Store = store
	return MiddlewareWithConfig(c)
}

// MiddlewareWithConfig returns a Sessions middleware with config.
// See `Middleware()`.
func MiddlewareWithConfig(config Config) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}
	if config.Store == nil {
		panic("echo: session middleware requires store")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			c.Set(key, config.Store)
			return next(c)
		}
	}
}
