package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/website-openid-proxy/internal/echosessions"
)

type Config struct {
	Skipper middleware.Skipper
}

func CheckAuthWithConfig(cfg Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}

			sess, err := echosessions.Get(loggedInCookieName, c)
			if err != nil {
				return c.Redirect(302, "/auth/login")
			}

			log.Ctx(c.Request().Context()).Info().Str("email", sess.Get("email")).Msg("user request")

			return next(c)
		}
	}
}
