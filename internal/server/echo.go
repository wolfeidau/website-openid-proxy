package server

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// LoginSkipper used to avoid running middleware for login requests
func LoginSkipper(prefix string) middleware.Skipper {
	return func(c echo.Context) bool {
		return strings.HasPrefix(c.Path(), prefix)
	}
}
